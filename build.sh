#!/bin/bash

set -e

CWD=`pwd`

build_mac64() {
    echo "Building mac64 components..."
    sh build_staticserver.sh M64
    sh build_htmlbuilder.sh M64

    echo "Creating zip archive..."
    cd $CWD/build
    cd Mac
    zip -qr ../mrlf-syllabus-viewer-Mac.zip .
    echo "Zip archive created"
    cd ..

    echo "Removing temporary build files"
    rm -r Mac

    echo "Build Completed."
}
build_windows64() {
    echo "Building Windows64 components..."
    # sh build_sms.sh
    sh build_staticserver.sh W64
    sh build_htmlbuilder.sh W64

    echo "Creating zip archive..."
    cd $CWD/build
    cd Win64
    zip -qr ../mrlf-syllabus-viewer-Win64.zip .
    echo "Zip archive created"
    cd ..

    echo "Removing temporary build files"
    rm -r Win64

    echo "Build Completed."
}
build_linux64() {
    echo "Building Linux64 components..."
    # sh build_sms.sh
    sh build_staticserver.sh L64
    sh build_htmlbuilder.sh L64

    echo "Creating zip archive..."
    cd $CWD/build
    cd Linux64
    zip -qr ../mrlf-syllabus-viewer-Linux64.zip .
    echo "Zip archive created"
    cd ..

    echo "Removing temporary build files"
    rm -r Linux64

    echo "Build Completed."
}
build_all() {
    echo "Building all components..."
    
    rm -rf $OUTPUT
    
    sh build_staticserver.sh
    sh build_htmlbuilder.sh

    echo "Creating zip archives..."
    cd $CWD/build

    cd Mac
    zip -qr ../mrlf-syllabus-viewer-Mac.zip .
    cd ..

    cd Win64
    zip -qr ../mrlf-syllabus-viewer-Win64.zip .
    cd ..

    cd Linux64
    zip -qr ../mrlf-syllabus-viewer-Linux64.zip .
    cd ..

    echo "Zip archives created"

    echo "Removing temporary build files"
    rm -r Mac Win64 Linux64

    echo "Build Completed."
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
    build_all
fi
