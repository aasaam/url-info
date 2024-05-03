package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/publicsuffix"
	"golang.org/x/text/language"
)

const singlePixelPNG = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z/C/HgAGgwJ/lK3Q6wAAAABJRU5ErkJggg=="

var stripTagger = bluemonday.StripTagsPolicy()

var temporaryPath = "/tmp"

var rtlLanguages = []string{"ar", "dv", "fa", "he", "ps", "ur", "yi"}
var rtlLanguagesMap map[string]bool

func init() {
	TEMPORARY_PATH := os.Getenv("ASM_URL_INFO_TEMPORARY_PATH")
	if TEMPORARY_PATH != "" {
		temporaryPath = TEMPORARY_PATH
	}
	rtlLanguagesMap = make(map[string]bool)
	for _, v := range rtlLanguages {
		rtlLanguagesMap[v] = true
	}
}

func validPublicSuffix(u *url.URL) bool {
	host := u.Hostname()
	eTLD, icann := publicsuffix.PublicSuffix(host)

	if icann {
		return true
	} else if strings.IndexByte(eTLD, '.') >= 0 {
		return true
	}

	return false
}

func urlWithoutPathAndQuery(u *url.URL) *url.URL {
	un, _ := url.Parse(u.String())

	un.RawQuery = ""
	un.Fragment = ""
	un.RawPath = ""
	un.Path = ""

	uns, _ := url.Parse(un.String())

	return uns
}

func execute(command string, arg ...string) ([]byte, int, error) {
	cmd := exec.Command(command, arg...)
	var stdOut bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, cmd.Process.Pid, errors.New(err.Error() + ": " + stderr.String())
	}
	return stdOut.Bytes(), cmd.Process.Pid, nil
}

func sanitizeLanguage(locale string) string {
	tag, err := language.Parse(locale)
	if err != nil {
		return ""
	}
	base, _ := tag.Base()
	return base.String()
}

func getLanguageDirection(lang string) string {
	if _, ok := rtlLanguagesMap[lang]; ok {
		return "rtl"
	}
	return "ltr"
}

func faviconConvert(u *url.URL) string {
	faviconURL := urlWithoutPathAndQuery(u)
	faviconURL.Path = "/favicon.ico"
	url := faviconURL.String()
	cacheFile := cachePath(url, "txt")
	if fileExist(cacheFile) {
		b, bErr := os.ReadFile(cacheFile)
		if bErr != nil {
			panic(bErr)
		}
		return string(b)
	}
	base64Image := singlePixelPNG
	cacheFileICO := cachePath(url, "ico")
	_, _, curlErr := execute("curl", "-Lsk", "-o", cacheFileICO, faviconURL.String())
	if curlErr == nil {
		cacheFilePNG := cachePath(url, "png")
		_, _, convertErr := execute("convert", cacheFileICO, "-thumbnail", "16x16", "-alpha", "on", "-flatten", "-strip", cacheFilePNG)
		if convertErr == nil {
			f, fErr := os.Open(cacheFilePNG)
			if fErr != nil {
				panic(fErr)
			}
			reader := bufio.NewReader(f)
			content, contentErr := io.ReadAll(reader)
			if contentErr != nil {
				panic(contentErr)
			}
			base64Image = base64.StdEncoding.EncodeToString(content)
		}
	}
	writeErr := os.WriteFile(cacheFile, []byte(base64Image), 0666)
	if writeErr != nil {
		panic(writeErr)
	}
	return base64Image
}

func imageConvert(u *url.URL, processParams urlInfoProcess) string {
	url := u.String()
	quality := processParams.ImageQuality
	if quality < 10 {
		quality = 10
	} else if quality > 95 {
		quality = 95
	}

	qualityString := strconv.Itoa(quality)

	cacheFile := cachePath(url+":"+processParams.ImageResize+":"+qualityString, "jpg")
	if fileExist(cacheFile) {
		return cacheFile
	}

	// download
	cacheFileImage := cachePath(url, "")
	if !fileExist(cacheFileImage) {
		_, _, curlErr := execute("curl", "-Lsk", "-o", cacheFileImage, url)
		if curlErr != nil {
			return ""
		}
	}

	// optimize
	_, _, convertErr := execute("convert", cacheFileImage, "-resize", processParams.ImageResize, "-sampling-factor", "4:2:0", "-quality", qualityString, "-interlace", "JPEG", "-colorspace", "sRGB", "-flatten", "-strip", cacheFile)
	if convertErr == nil && fileExist(cacheFile) {
		return cacheFile
	}

	return ""
}

func hash(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func sanitizeString(str string) string {
	return strings.TrimSpace(str)
}

func sanitizeText(t string) string {
	return sanitizeString(stripTagger.Sanitize(t))
}

func sanitizeIntPointer(s string) *int {
	i, e := strconv.Atoi(s)
	if e == nil && i > 0 {
		return &i
	}
	return nil
}

func parseKeywords(inpKeywords string) []string {
	r := []string{}
	if inpKeywords == "" {
		return r
	}

	ks := strings.Split(inpKeywords, ",")
	for _, k := range ks {
		v := sanitizeText(k)
		if len(v) >= 1 {
			r = append(r, v)
		}
	}

	if len(r) > 10 {
		return r[0:10]
	}

	return r
}

func cachePath(s string, suffix string) string {
	return fmt.Sprintf("%s/%s.%s", temporaryPath, hash(s), suffix)
}

func fileExist(p string) bool {
	if st, err := os.Stat(p); err == nil {
		return !st.IsDir()
	}
	return false
}
