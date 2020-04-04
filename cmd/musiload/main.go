package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/prateek2211/musiload-go"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: musiload <url>")
		return
	}
	urlIp := os.Args[1]
	u, _ := url.Parse(urlIp)

	downloader := musiload.NewDownloader(*u)
	err := downloader.Download()
	if err != nil {
		log.Println(err.Error())
	}

}
