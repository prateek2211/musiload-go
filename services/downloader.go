package services

import (
	"github.com/golang/groupcache/lru"
	"github.com/grafov/m3u8"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func ParsePlaylist(mediaPlaylistUrl string, ts chan string) {
	cache := lru.New(1024)
	ipUrl, _ := url.Parse(mediaPlaylistUrl)
	for {
		response, err := http.Get(mediaPlaylistUrl)
		if err != nil {
			log.Println(err.Error())
			time.Sleep(time.Duration(3) * time.Second)
		}
		playlist, listtype, err := m3u8.DecodeFrom(response.Body, true)
		if err != nil {
			log.Fatal(err.Error())
		}
		if listtype == m3u8.MEDIA {
			mp := playlist.(*m3u8.MediaPlaylist)
			for _, segment := range (mp.Segments) {
				var tsURI string
				if segment == nil {
					break;
				}
				if strings.HasPrefix(segment.URI, "http") {
					tsURI = segment.URI
				} else {
					segUrl, err := ipUrl.Parse(segment.URI)
					if err != nil {
						log.Fatal(err.Error())
					}
					tsURI, _ = url.QueryUnescape(segUrl.String())
				}
				_, hit := cache.Get(tsURI)
				if !hit {
					cache.Add(tsURI, nil)
					ts <- tsURI
				}
			}
			if mp.Closed {
				close(ts)
				return
			} else {
				time.Sleep(time.Duration(int64(mp.TargetDuration * 1000000000)))
			}
		}
	}
}
func DownloadTS(ts chan string, fileName string) {
	file, _ := os.Create(fileName)
	for stream := range (ts) {
		resp, err := http.Get(stream)
		if err != nil {
			log.Fatal(err.Error())
		}
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
