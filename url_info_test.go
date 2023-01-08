package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewURL000(t *testing.T) {
	u, uE := newURL("http://parsnews.com/fa/tiny/news-684283")
	if uE != nil {
		t.Error(uE)
		return
	}
	u.process(urlInfoProcess{
		ImageResize:  "720x240!",
		ImageQuality: 30,
	})
	b, _ := json.Marshal(u)
	fmt.Println(string(b))
}
