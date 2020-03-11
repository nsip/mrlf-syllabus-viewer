#!/bin/bash

set -e

CWD=`pwd`


do_build() {
	echo "Building static webserver..."
    # rm -rf $OUTPUT
	mkdir -p $OUTPUT
	cd $CWD
	cd ./cmd/staticserver
	go get
	GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUTPUT/$HARNESS
	cd ..
}

do_zip() {
	cd $OUTPUT
	cd ..
	zip -qr ../$ZIP staticserver
	cd $CWD
}

build_mac64() {
	# MAC OS X (64 only)
	echo "Building Mac binaries..."
	GOOS=darwin
	GOARCH=amd64
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Mac
	HARNESS=staticserver
	ZIP=staticserver-Mac.zip
	do_build
	echo "...all Mac binaries built..."
}


build_windows64() {
	# WINDOWS 64
	echo "Building Windows64 binaries..."
	GOOS=windows
	GOARCH=amd64
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Win64
	HARNESS=staticserver.exe
	ZIP=staticserver-Win64.zip
	do_build
	echo "...all Windows64 binaries built..."
}

build_linux64() {
	# LINUX 64
	echo "Building Linux64 binaries..."
	GOOS=linux
	GOARCH=amd64
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Linux64
	HARNESS=staticserver
	ZIP=staticserver-Linux64.zip
	do_build
	echo "...all Linux64 binaries built..."
}


if [ "$1" = "L64"  ]
then
    build_linux64
elif [ "$1" = "W64"  ]
then
    build_windows64
elif [ "$1" = "M64"  ]
then
    build_mac64
else
    build_mac64
    build_windows64
    build_linux64
fi
