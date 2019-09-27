package hungama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/grafov/m3u8"
	"github.com/prateek2211/musiload-go/platforms/saavn"
	"github.com/prateek2211/musiload-go/services"
	"log"
	"net/url"
	"strings"
)

type song struct {
	fileUrl   string
	mediaUrl  string
	songTitle string
}

func ParseAndDownload(link string) {
	splits := strings.Split(link, "/")
	uid := splits[len(splits)-2]
	pu, _ := url.Parse(link)
	var isAlbum = false
	if strings.Split(pu.Path, "/")[1] == "album" {
		isAlbum = true
	}
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:69.0) Gecko/20100101 Firefox/69.0"
	var songs []song
	var index = 0
	c.OnResponse(func(r *colly.Response) {
		if strings.Contains(r.Request.URL.String(), "/audio-player-data/") {
			var data []interface{}
			err := json.Unmarshal(r.Body, &data)
			if err != nil {
				log.Println(err.Error())
			}
			for x, _ := range data {
				var s song
				s.fileUrl = data[x].(map[string]interface{})["file"].(string)
				s.songTitle = data[x].(map[string]interface{})["song_name"].(string)
				songs = append(songs, s)
			}
		}
		if strings.Contains(r.Request.URL.String(), "/mdnurl/song") {
			var data map[string]interface{}
			err := json.Unmarshal(r.Body, &data)
			if err != nil {
				log.Println(err.Error())
			}
			songs[index].mediaUrl = data["response"].(map[string]interface{})["media_url"].(string)
		}
		if r.Headers.Get("content-type") == "application/vnd.apple.mpegurl" {
			fmt.Println("Downloading ...")
			playist, _, _ := m3u8.DecodeFrom(bytes.NewReader(r.Body), true)
			mp := playist.(*m3u8.MasterPlaylist)
			ts := make(chan string, 1024)
			go services.ParsePlaylist(mp.Variants[len(mp.Variants)-1].URI, ts)
			services.DownloadTS(ts, songs[index].songTitle+".mp3")
		}
		if r.Headers.Get("content-type") == "audio/mpeg" {
			saavn.DownloadAudio(songs[index].mediaUrl, songs[index].songTitle+".mp3")
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	err := c.Post("https://www.hungama.com/user/loggedin-status", nil)
	if err != nil {
		log.Println(err.Error())
	}
	if isAlbum {
		err = c.Visit(fmt.Sprintf("https://www.hungama.com/audio-player-data/%s/%s?_country=IN", "album", uid))
	} else {
		err = c.Visit(fmt.Sprintf("https://www.hungama.com/audio-player-data/%s/%s?_country=IN", "track", uid))
	}
	u := url.URL{}
	u.Host = "https://ping.hungama.com/t.js"
	u.Query().Add("_hn", "hungama_web")
	u.Query().Add("_url", link)
	err = c.Visit(u.EscapedPath())
	for _, s := range songs {
		err = c.Visit(s.fileUrl)
		err = c.Visit(songs[index].mediaUrl)
		index++
	}
}
