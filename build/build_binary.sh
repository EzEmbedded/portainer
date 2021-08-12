#!/bin/sh

binary="portainer"
mkdir -p dist

cd 'api/cmd/portainer'

go get -u -t -d -v ./...
GOOS=$1 GOARCH=$2 CGO_ENABLED=0 GO111MODULE=on go build -a --installsuffix cgo --ldflags '-s'

mv "${binary}" "../../../dist/portainer"