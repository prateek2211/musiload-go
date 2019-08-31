package gaana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/grafov/m3u8"
	"github.com/prateek2211/musiload-go/services"
	"github.com/prateek2211/musiload-go/utils"
	"log"
)

func ParseAndDownload(songUrl string) {
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
		decodedURL = utils.StripCtlAndExtFromUnicode(decodedURL)
		err = c.Visit(decodedURL)
		if err != nil {
			log.Fatal(err.Error())
		}
	})

	c.OnResponse(func(response *colly.Response) {
		if (response.Headers.Get("content-type") == "application/vnd.apple.mpegurl") {
			fmt.Println("Downloading ...")
			playist, _, _ := m3u8.DecodeFrom(bytes.NewReader(response.Body), true)
			mp := playist.(*m3u8.MasterPlaylist)
			ts := make(chan string, 1024)
			go services.ParsePlaylist(mp.Variants[0].URI, ts)
			services.DownloadTS(ts, songTitle+".mp3")
			//err := ioutil.WriteFile(songTitle+".m3u8", response.Body, 0644)
			//if err != nil {
			//	log.Fatal(err.Error())
			//}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	err := c.Visit(songUrl)
	if err != nil {
		log.Fatal(err.Error())
	}
}
