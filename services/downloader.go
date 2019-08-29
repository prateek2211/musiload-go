package services

import (
	"github.com/grafov/m3u8"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func ParsePlaylist(mediaPlaylistUrl string, ts chan string) {
	ipUrl, _ := url.Parse(mediaPlaylistUrl)
	response, err := http.Get(mediaPlaylistUrl)
	if err != nil {
		log.Fatal(err.Error())
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
			ts <- tsURI
		}
		close(ts)
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
