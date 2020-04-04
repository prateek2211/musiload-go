#!/bin/sh

if [[ $EUID > 0 ]]; then 
  echo "The installation script needs to be run as root."
  exit 1
else
  mkdir bin
  /usr/local/go/bin/go build -o bin -ldflags "-s -w" ./... > /dev/null
  cp ./bin/musiload /usr/local/bin
  chmod 755 /usr/local/bin/musiload
fi
