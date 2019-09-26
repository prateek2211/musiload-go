package saavn

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/inhies/go-bytesize"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func ParseAndDownload(songUrl string) {
	var songTitle string
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:68.0) Gecko/20100101 Firefox/68.0"
	c.OnHTML(".meta-info", func(e *colly.HTMLElement) {
		songInfoText := e.ChildText("div")
		var songInfo map[string]interface{}
		err := json.Unmarshal([]byte(songInfoText), &songInfo)
		if err != nil {
			log.Fatal(err.Error())
		}
		songUrl := songInfo["url"].(string)
		songTitle = songInfo["title"].(string)
		//Sets the cookie
		err = c.Visit("https://www.jiosaavn.com/stats.php?ev=site:browser:fp&fp=a49a8a770fe85d2031bbf09d23be7e86")
		if err != nil {
			log.Fatal(err.Error())
		}
		err = c.Post("https://www.jiosaavn.com/api.php", map[string]string{"url": songUrl, "__call": "song.generateAuthToken", "_marker": "false", "_format": "json", "bitrate": "128"})
		if err != nil {
			log.Fatal(err.Error())
		}
	})
	c.OnResponse(func(r *colly.Response) {
		if r.Request.URL.String() == "https://www.jiosaavn.com/api.php" {
			var response map[string]interface{}
			err := json.Unmarshal(r.Body, &response)
			if err != nil {
				log.Fatal(err.Error())
			}
			DownloadAudio(response["auth_url"].(string), songTitle+".mp3")
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

type Sizecounter struct {
	size uint64
}

func (c *Sizecounter) Write(d []byte) (int, error) {
	c.size += uint64(len(d))
	printDownloadedSize(c.size)
	return len(d), nil
}

func printDownloadedSize(u uint64) {
	size := bytesize.New(float64(u)).String()
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rDownloading... %s complete", size)
}

func DownloadAudio(url string, fileName string) {
	response, err := http.Get(url)
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = io.Copy(file, io.TeeReader(response.Body, &Sizecounter{}))
	if err != nil {
		log.Fatal(err.Error())
	}
	println()
}
