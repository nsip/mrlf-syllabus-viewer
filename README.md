# mrlf-syllabus-viewer
creates generated views of machine-readable syllabus

## Building from source
All source code available using standard go get:

```go get github.com/nsip/mrlf-syllabus-viewer```

to generate binaries for your platform run 

```build.sh```

in the root directory.
(Build script needs a unix-like environment)

## Binary packages
Pre-built binaries are avalable [here](https://github.com/nsip/mrlf-syllabus-viewer/releases/latest)

## Installing
Unpack the zip file for your platform to a suitable folder.

The package contains two binaries and various directories of support files.

The htmlbuilder(.exe) component generates the html representations of the json syllabus.

The statiserver(.exe) component is a simple webserver that will serve the generated html files, supporting css & 
javascript so you can review the output. The staticserver runs on port 1323.

When htmlbuilder runs it generates all resources under the ```/public```folder. These are now static presentaitons of the
syllabus.
If you choose to run staticserver it simply serves whatever it finds in the ```/public``` folder.
You can equally copy the contents of the public folder to be served by any server of your choice.

To generate and serve the html views of the syllabus:

unpack the zip to a suitable folder,

you’ll need to open two console sessions at the folder you unzipped to…

in session one launch the staticserver   - just a simple web server so you can see the published reports.

```./staticserver(.exe)```

in the other run the html builder

```./htmlbuilder(.exe)```

or 

```./htmlbuilder(.exe) -css```

if you want to re-generate the contextual css used by the content view (for exmaple if the structure of your 
input data has changed).

this will take any JSON syllabus in the /input folder and render the three views, the nesa2.json and output.json files are 
bundled already.

assuming everything works, you can then point your browser to the following reports:

http://localhost:1323/nesa2/audit.html 

http://localhost:1323/output/audit.html 

- ‘audit’ view, shows whole structure and allows you to toggle visibility of class nestings 

http://localhost:1323/nesa2/glossary.html

http://localhost:1323/output/glossary.html

- glossary view, extracts the glossary and makes it searchable 

http://localhost:1323/nesa2/content.html

http://localhost:1323/output/content.html

- content view shows selected items of content (based on css)




