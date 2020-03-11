#!/bin/bash

set -e

CWD=`pwd`

do_build() {
	echo "Building htmlbuilder..."
	mkdir -p $OUTPUT
	cd $CWD
	cd ./cmd/htmlbuilder
	GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUTPUT/$HARNESS
	cd ..
	rsync -a ../input ../resources ../templates $OUTPUT/
	# rsync -a napval/students.csv $OUTPUT/
}

do_zip() {
	cd $OUTPUT
	cd ..
	#zip -qr ../$ZIP napval
	cd $CWD
}

build_mac64() {
	# MAC OS X (64 only)
	echo "Building Mac binaries..."
	GOOS=darwin
	GOARCH=amd64
	LDFLAGS="-s -w"
    OUTPUT=$CWD/build/Mac
	HARNESS=htmlbuilder
	ZIP=htmlbuilder-Mac.zip
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
	HARNESS=htmlbuilder.exe
	ZIP=htmlbuilder-Win64.zip
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
	HARNESS=htmlbuilder
	ZIP=htmlbuilder-Linux64.zip
	do_build
	echo "...all Linux64 binaries built..."
}

# TODO ARM
# GOOS=linux GOARCH=arm GOARM=7 go build -o $CWD/build/LinuxArm7/go-nias/aggregator

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
