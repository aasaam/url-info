package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewURL000(t *testing.T) {

	urls := []string{
		"https://raw.githubusercontent.com/niallkennedy/open-graph-protocol-examples/master/article.html",
		"https://raw.githubusercontent.com/niallkennedy/open-graph-protocol-examples/master/audio.html",
		"https://raw.githubusercontent.com/niallkennedy/open-graph-protocol-examples/master/video-movie.html",
	}

	for _, ur := range urls {
		u, uE := newURL(ur)
		if uE != nil {
			t.Error(uE)
			return
		}
		u.process(urlInfoProcess{
			URL:          ur,
			ImageResize:  "720x240!",
			ImageQuality: 30,
		})
		b, _ := json.MarshalIndent(u, "", "  ")
		fmt.Println(string(b))
	}

}
