package main

import (
	"bytes"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/PuerkitoBio/purell"
)

const maxHTMLSize = 1024 * 1024 * 5

var defaultURLSanitize purell.NormalizationFlags = purell.FlagLowercaseScheme |
	purell.FlagLowercaseHost |
	purell.FlagDecodeUnnecessaryEscapes |
	purell.FlagEncodeNecessaryEscapes |
	purell.FlagRemoveDefaultPort |
	purell.FlagRemoveEmptyQuerySeparator |
	purell.FlagRemoveDotSegments |
	purell.FlagRemoveDuplicateSlashes |
	purell.FlagSortQuery |
	purell.FlagRemoveUnnecessaryHostDots |
	purell.FlagRemoveEmptyPortSeparator |
	purell.FlagUppercaseEscapes

type urlInfo struct {
	url        *url.URL
	short_link *url.URL
	image_url  *url.URL
	canonical  *url.URL

	Title            string   `json:"title,omitempty"`
	Language         string   `json:"lang,omitempty"`
	Direction        string   `json:"dir,omitempty"`
	Keywords         []string `json:"keywords,omitempty"`
	Description      string   `json:"description,omitempty"`
	IconPath         string   `json:"icon_path,omitempty"`
	ImagePath        string   `json:"image_path,omitempty"`
	URL              string   `json:"url,omitempty"`
	Icon             string   `json:"icon,omitempty"`
	ShortLink        string   `json:"short_link,omitempty"`
	ImageURL         string   `json:"image_url,omitempty"`
	ImageConvertPath string   `json:"image_optimized,omitempty"`
	CanonicalURL     string   `json:"canonical,omitempty"`
}

type urlInfoProcess struct {
	URL          string `json:"url"`
	ImageResize  string `json:"image_resize"`
	ImageQuality int    `json:"image_quality"`
}

func (uip *urlInfoProcess) cache() string {
	qualityString := strconv.Itoa(uip.ImageQuality)
	params := []string{uip.URL, uip.ImageResize, qualityString}
	return hash(strings.Join(params[:], ","))
}

func newURL(u string) (*urlInfo, error) {
	normalizeURLString, normalizeURLStringErr := purell.NormalizeURLString(u, defaultURLSanitize)
	if normalizeURLStringErr != nil {
		return nil, normalizeURLStringErr
	}

	ur, urE := url.Parse(normalizeURLString)

	if urE != nil {
		return nil, urE
	}

	if !validPublicSuffix(ur) {
		return nil, errors.New("invalid public suffix")
	}

	o := urlInfo{
		url: ur,
	}

	return &o, nil
}

func (ui *urlInfo) process(processInfo urlInfoProcess) error {
	client := http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
		},
	}

	resp, respErr := client.Get(ui.url.String())

	if respErr != nil {
		return respErr
	}

	ui.url = resp.Request.URL
	ui.URL = resp.Request.URL.String()

	limitedReader := &io.LimitedReader{R: resp.Body, N: maxHTMLSize}
	body, bodyErr := io.ReadAll(limitedReader)
	if bodyErr != nil {
		return bodyErr
	}

	doc, docErr := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if docErr != nil {
		return docErr
	}

	// lang
	lang, langExist := doc.Find("html").First().Attr("lang")
	if langExist {
		ui.Language = sanitizeLanguage(lang)
	}

	ui.Direction = getLanguageDirection(ui.Language)

	// title
	titleOG, titleOGExist := doc.Find("head meta[property='og:title'],head meta[name='og:title']").First().Attr("content")
	if titleOGExist {
		ui.Title = sanitizeText(titleOG)
	} else {
		ui.Title = sanitizeText(doc.Find("head title").First().Text())
	}

	// description
	descriptionOG, descriptionOGExist := doc.Find("head meta[property='og:description'],head meta[name='og:description']").First().Attr("content")
	if descriptionOGExist {
		ui.Description = sanitizeText(descriptionOG)
	} else {
		description, descriptionExist := doc.Find("head meta[name='description']").First().Attr("content")
		if descriptionExist {
			ui.Description = sanitizeString(description)
		}
	}

	keywords, keywordsExist := doc.Find("head meta[name='keywords']").First().Attr("content")
	if keywordsExist {
		ui.Keywords = parseKeywords(keywords)
	}

	canonical, canonicalExist := doc.Find("head link[rel='canonical']").First().Attr("href")
	if canonicalExist {
		canonicalURLString, canonicalURLStringErr := purell.NormalizeURLString(canonical, defaultURLSanitize)
		if canonicalURLStringErr == nil {
			ui.canonical, _ = url.Parse(canonicalURLString)
			ui.CanonicalURL = ui.canonical.String()
		}
	}

	shortLink, shortLinkExist := doc.Find("head link[rel='shortlink']").First().Attr("href")
	if shortLinkExist {
		shortLinkURLString, shortLinkURLStringErr := purell.NormalizeURLString(shortLink, defaultURLSanitize)
		if shortLinkURLStringErr == nil {
			ui.short_link, _ = url.Parse(shortLinkURLString)
			ui.ShortLink = ui.short_link.String()
		}
	}

	ui.Icon = "data:image/png;base64," + faviconConvert(ui.url)

	imageOg, imageOgExist := doc.Find("head meta[property='og:image'],head meta[name='og:image']").First().Attr("content")
	if imageOgExist {
		imageOgURLString, imageOgURLStringErr := purell.NormalizeURLString(imageOg, defaultURLSanitize)
		if imageOgURLStringErr == nil {
			ui.image_url, _ = url.Parse(imageOgURLString)
			ui.ImageURL = ui.image_url.String()

			ui.ImageConvertPath = imageConvert(ui.image_url, processInfo)
		}
	}

	return nil
}
