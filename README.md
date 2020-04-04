# Musiload [![Build Status](https://travis-ci.com/prateek2211/musiload-go.svg?branch=master)](https://travis-ci.com/prateek2211/musiload-go)

Musiload downloads songs/playlists from famous websites like [gaana](gaana.com), [hungama](www.hungama.com), [saavn](www.jiosaavn.com)

Just enter the url of song/playlist and musiload does the job for you. 


## Usage

* ### Install directly from release:
    * Grab the binary suitable for your OS [here](https://github.com/prateek2211/musiload-go/releases)

    * Make the binary executable
```shell script
$ chmod +x musiload
```
* ### Build from source:

<b> Prerequisite </b>

* Golang (1.11 or later). You can see the instructions [here](https://golang.org/dl/) to download

```bash
$ git clone https://github.com/prateek2211/musiload-go.git
$ mkdir bin
$ go build -o bin ./...
```

Run program and enter the website url of the song

```bash
$ .bin/musiload <URL>
```

The song will be downloaded in the Music directory


## Todo

* Add functionality for other famous websites
* Convert the program to a web app
* Add feature to search music from websites
