package main

import (
	"fmt"
	"github.com/prateek2211/musiload-go/platforms/gaana"
	"github.com/prateek2211/musiload-go/platforms/saavn"
	"log"
	"net/url"
)

func main() {
	fmt.Println("Enter the website url:")
	var urlIp string
	_, err := fmt.Scanln(&urlIp)
	if err != nil {
		log.Fatal(err.Error())
	}
	u, _ := url.Parse(urlIp)
	if u.Hostname() == "gaana.com" {
		gaana.ParseAndDownload(urlIp)
	}
	if u.Hostname() == "www.jiosaavn.com" {
		saavn.ParseAndDownload(urlIp)
	}
}
