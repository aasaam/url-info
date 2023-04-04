package main

import (
	"fmt"
	"net/url"
	"testing"
)

func TestUrlWithoutPathAndQuery(t *testing.T) {
	u, _ := url.Parse("http://www.test.com/url?foo=bar&foo=baz#this_is_fragment")
	fmt.Println(urlWithoutPathAndQuery(u))
	fmt.Println(urlWithoutPathAndQuery(u))
}

func TestUrlValidate(t *testing.T) {

	validURLs := []string{
		"https://google.com/path/file.ext?foo=bar#a=b",
		"https://www.books.amazon.co.uk/path/file.ext?foo=bar#a=b",
	}

	inValidURLs := []string{
		"https://google.com-not-exist/path/file.ext?foo=bar#a=b",
		"http://127.0.0.1/path/file.ext?foo=bar#a=b",
	}

	for _, str := range validURLs {
		u, _ := url.Parse(str)

		if !validPublicSuffix(u) {
			t.Errorf("must valid")
		}
	}

	for _, str := range inValidURLs {
		u, _ := url.Parse(str)

		if validPublicSuffix(u) {
			t.Errorf("must invalid")
		}
	}
}
