package musiload

import (
	"github.com/prateek2211/musiload-go/platforms/gaana"
	"github.com/prateek2211/musiload-go/platforms/hungama"
	"github.com/prateek2211/musiload-go/platforms/saavn"
	"net/url"
)

type Downloader interface {
	Download() error
}

const HostnameHungama = "www.hungama.com"

const HostnameGaana = "gaana.com"

const HostnameSaavn = "www.jiosaavn.com"

func NewDownloader(url url.URL) Downloader {
	switch url.Hostname() {
	case HostnameSaavn:
		return saavn.SaavnDownloader{url}
	case HostnameGaana:
		return gaana.GaanaDownloader{url}
	case HostnameHungama:
		return hungama.HungamaDownloader{url}
	default:
		return nil
	}
}
