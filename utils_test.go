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
