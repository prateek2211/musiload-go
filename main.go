package main

import (
	"fmt"
	"github.com/prateek2211/musiload-go/platforms/gaana"
	"github.com/prateek2211/musiload-go/platforms/hungama"
	"github.com/prateek2211/musiload-go/platforms/saavn"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: musiload <url>")
		return
	}
	urlIp := os.Args[1]
	//fmt.Println("Enter the website url:")
	//var urlIp string
	//_, err := fmt.Scanln(&urlIp)
	//if err != nil {
	//	log.Fatal(err.Error())
	//}
	u, _ := url.Parse(urlIp)
	if u.Hostname() == "gaana.com" {
		gaana.ParseAndDownload(urlIp)
	}
	if u.Hostname() == "www.jiosaavn.com" {
		saavn.ParseAndDownload(urlIp)
	}
	if u.Hostname() == "www.hungama.com" {
		hungama.ParseAndDownload(urlIp)
	}
}
