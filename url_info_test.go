package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewURL000(t *testing.T) {
	u, uE := newURL("http://news.yahoo.com/")
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
