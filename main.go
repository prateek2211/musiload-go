package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/prateek2211/musiload-go/utils"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"io/ioutil"
	"log"
)

func main() {
	var songTitle string
	c := colly.NewCollector()

	c.OnHTML("ul.s_l.artworkload", func(e *colly.HTMLElement) {
		link := e.ChildText("span[data-type=\"playSong\"]")
		var data map[string]interface{}
		err := json.Unmarshal([]byte(link), &data)
		if err != nil {
			log.Fatal(err.Error())
		}
		songTitle = data["title"].(string)
		highQuality := data["path"].(map[string]interface{})["high"].([]interface{})[0].(map[string]interface{})
		cip := highQuality["message"].(string)
		decodedURL := utils.Decrypt([]byte(cip), "g@1n!(f1#r.0$)&%", "asd!@#!@#@!12312")
		decodedURL = stripCtlAndExtFromUnicode(decodedURL)
		err = c.Visit(decodedURL)
		if err != nil {
			log.Fatal(err.Error())
		}
	})

	c.OnResponse(func(response *colly.Response) {
		if (response.Headers.Get("content-type") == "application/vnd.apple.mpegurl") {
			fmt.Println("Hello")
			err := ioutil.WriteFile(songTitle+".m3u8", response.Body, 0644)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	var url string
	_, err := fmt.Scanln(&url)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = c.Visit(url)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func stripCtlAndExtFromUnicode(str string) string {
	isOk := func(r rune) bool {
		return r < 32 || r >= 127
	}
	// The isOk filter is such that there is no need to chain to norm.NFC
	t := transform.Chain(norm.NFKD, transform.RemoveFunc(isOk))
	// This Transformer could also trivially be applied as an io.Reader
	// or io.Writer filter to automatically do such filtering when reading
	// or writing data anywhere.
	str, _, _ = transform.String(t, str)
	return str
}
