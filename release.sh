
cd tools; go build release.go; cd ..
./tools/release mrlf-syllabus-viewer > version/version.go
sh build.sh
./tools/release mrlf-syllabus-viewer mrlf-syllabus-viewer-Mac.zip build/mrlf-syllabus-viewer-Mac.zip
./tools/release mrlf-syllabus-viewer mrlf-syllabus-viewer-Win64.zip build/mrlf-syllabus-viewer-Win64.zip
./tools/release mrlf-syllabus-viewer mrlf-syllabus-viewer-Linux64.zip build/mrlf-syllabus-viewer-Linux64.zip
