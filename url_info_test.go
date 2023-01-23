package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewURL000(t *testing.T) {
	u, uE := newURL("https://www.eghtesadnews.com/%D8%A8%D8%AE%D8%B4-%D8%A7%D8%AE%D8%A8%D8%A7%D8%B1-%D8%B3%DB%8C%D8%A7%D8%B3%DB%8C-57/550397-%D9%88%D8%A7%DA%A9%D9%86%D8%B4-%D9%88%D8%B2%DB%8C%D8%B1-%DA%A9%D8%B4%D9%88%D8%B1-%D8%A8%D9%87-%D8%A7%D9%81%D8%B2%D8%A7%DB%8C%D8%B4-%D9%82%DB%8C%D9%85%D8%AA-%D8%AF%D9%84%D8%A7%D8%B1-%D9%81%D8%B1%D9%88%D9%BE%D8%A7%D8%B4%DB%8C-%D8%A2%D9%85%D8%B1%DB%8C%DA%A9%D8%A7-%D8%B3%D9%82%D9%88%D8%B7-%DA%A9%D8%A7%D8%AE-%D8%AF%D9%84%D8%A7%D8%B1-%D8%B1%D8%A7-%D8%AE%D9%88%D8%A7%D9%87%DB%8C%D9%85-%D8%AF%DB%8C%D8%AF")
	if uE != nil {
		t.Error(uE)
		return
	}
	u.process(urlInfoProcess{
		URL:          "https://www.eghtesadnews.com/%D8%A8%D8%AE%D8%B4-%D8%A7%D8%AE%D8%A8%D8%A7%D8%B1-%D8%B3%DB%8C%D8%A7%D8%B3%DB%8C-57/550397-%D9%88%D8%A7%DA%A9%D9%86%D8%B4-%D9%88%D8%B2%DB%8C%D8%B1-%DA%A9%D8%B4%D9%88%D8%B1-%D8%A8%D9%87-%D8%A7%D9%81%D8%B2%D8%A7%DB%8C%D8%B4-%D9%82%DB%8C%D9%85%D8%AA-%D8%AF%D9%84%D8%A7%D8%B1-%D9%81%D8%B1%D9%88%D9%BE%D8%A7%D8%B4%DB%8C-%D8%A2%D9%85%D8%B1%DB%8C%DA%A9%D8%A7-%D8%B3%D9%82%D9%88%D8%B7-%DA%A9%D8%A7%D8%AE-%D8%AF%D9%84%D8%A7%D8%B1-%D8%B1%D8%A7-%D8%AE%D9%88%D8%A7%D9%87%DB%8C%D9%85-%D8%AF%DB%8C%D8%AF",
		ImageResize:  "720x240!",
		ImageQuality: 30,
	})
	b, _ := json.Marshal(u)
	fmt.Println(string(b))
}
