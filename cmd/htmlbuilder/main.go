package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

//
// struct to hold all necessary
// meta-data such as paths, root name etc.
// derived from an input filename
//
type FileMeta struct {
	RootName string
	FileName string
	HtmlPath string
	CssPath  string
	JsPath   string
}

type TemplateData struct {
	FM   FileMeta
	JSON interface{}
}

//
// ***********************************************
// helper functions for templates
//
var isObjectArray = func(val interface{}) bool {
	switch val.(type) {
	case []map[string]interface{}:
		return true
	default:
		return false
	}
}

var isLeaf = func(val interface{}) bool {
	switch val.(type) {
	case bool, float64, string:
		return true
	default:
		return false
	}
}

var isArray = func(val interface{}) bool {
	switch val.(type) {
	case []interface{}:
		return true
	default:
		return false
	}
}

var isObject = func(val interface{}) bool {
	switch val.(type) {
	case map[string]interface{}:
		return true
	default:
		return false
	}
}

//
// **************************************************

var baseOutputFolder string
var resourcesFolder string

func main() {

	inputFolder := flag.String("i", "./input", "json input folder name")
	outputFolder := flag.String("o", "./public/", "root output folder name")

	//
	// flag to say don't regenerate if you already have something working
	//
	css := flag.Bool("css", false, "generate content-view css template")

	// TODO:
	// delete public folders & rebuild - flag!
	//

	flag.Parse()

	baseOutputFolder = *outputFolder
	resourcesFolder = "./resources"

	srcImageFolder := fmt.Sprintf("%s/%s", resourcesFolder, "image")
	destImgFolder := fmt.Sprintf("%s%s", baseOutputFolder, "images")
	err := os.MkdirAll(destImgFolder, os.ModePerm)
	if err != nil {
		fmt.Println("unable to create target image directory: ", err)
	}

	cerr := copyDirectory(srcImageFolder, destImgFolder)
	if cerr != nil {
		fmt.Println("unable to copy image directory: ", err)
	}

	// open folder and parse all json filenames
	// make sure default input folder exists
	err = os.MkdirAll(*inputFolder, os.ModePerm)
	if err != nil {
		log.Fatal("unable to read input directory: ", err)
	}

	fileNames, err := parseInputFolder(*inputFolder)
	if err != nil {
		log.Fatal("unable to read json files from input folder: ", err)
	}

	log.Println("files found: ", fileNames)

	fmList := make([]FileMeta, 0) // keep a lis tof all gneratd files to build outer nav.
	for _, fileName := range fileNames {

		//
		// TODO: error hndling, return-fatal
		// go-routine multi files & views
		// pool?
		//

		// create file meta-data
		fm := buildFileMeta(fileName)
		log.Println("fm: ", fm)
		// make output folders for each read file
		err := makeOutputFolders(fm)
		if err != nil {
			log.Println("unable to create output folders for: ", fileName, err)
		}

		//
		// read the json file
		//
		jsonBytes, err := ioutil.ReadFile(fm.FileName)
		if err != nil {
			log.Println("cannot read json file: "+fm.FileName, err)
		}

		var data interface{}

		if err := json.Unmarshal(jsonBytes, &data); err != nil {
			log.Println("cannot unmarshal json file: ", err)
		}

		//
		// gen. audit view
		//
		err = createAuditView(fm, data)
		if err != nil {
			log.Println("unable to create audit view: ", err)
		}

		//
		// extract glossary; try at syllabus level first, otherwise
		// get from KLA-level
		//
		var gData interface{}
		syllabusGlossary := gjson.GetBytes(jsonBytes, "children.0.glossary")
		if syllabusGlossary.Exists() {
			gData = syllabusGlossary.Value()
		} else {
			log.Println("unable to find a syllabus glossary; view will be empty.")
		}

		// gem. glossary view
		// nb. is only one! at either level...
		err = createGlossaryView(fm, gData)
		if err != nil {
			log.Println("unable to create glossary view: ", err)
		}

		//
		// remove glossaries from data, and re-parse
		//
		noGlossary, err := sjson.DeleteBytes(jsonBytes, "children.0.glossary")
		if err != nil {
			log.Println("unable to delete syllabus-level glossary: ", err)
		}

		if err := json.Unmarshal(noGlossary, &data); err != nil {
			log.Println("cannot unmarshal post-gloss. remove json file: ", err)
		}

		//
		// gen. content view
		//
		err = createContentView(fm, data, *css)
		if err != nil {
			log.Println("unable to create content view: ", err)
		}

		//
		// copy jquery support resources
		//
		srcFile := fmt.Sprintf("%s/js/jquery-3.3.1.min.js", resourcesFolder)
		dstFile := fmt.Sprintf("%s/%s/jquery-3.3.1.min.js", baseOutputFolder, fm.JsPath)
		err = copyFileContents(srcFile, dstFile)
		if err != nil {
			log.Println("unable to copy jquery support files: ", err)
		}
		//
		// minimap navigator
		//
		srcFile = fmt.Sprintf("%s/js/pagemap.min.js", resourcesFolder)
		dstFile = fmt.Sprintf("%s/%s/pagemap.min.js", baseOutputFolder, fm.JsPath)
		err = copyFileContents(srcFile, dstFile)
		if err != nil {
			log.Println("unable to copy pagemap support files: ", err)
		}

		fmList = append(fmList, fm)

	}

	// mini-map
	// minimap css

	// progress bars

	// outer nav index.html
}

//
// creates a more human-readable summary of the
// content with only some textual elements rendered.
//
func createContentView(fm FileMeta, data interface{}, genCss bool) error {

	// create the output file
	outFileName := fmt.Sprintf("%s%s/%s", baseOutputFolder, fm.RootName, "content.html")
	log.Println("content-view output file: ", outFileName)
	html, err := os.Create(outFileName)
	if err != nil {
		log.Fatal("unable to create content-view html file: ", outFileName)
	}
	defer html.Close()

	// also create contextual css file if requested
	var css *os.File
	if genCss {
		cssFileName := resourcesFolder + "/css/content.css"
		log.Println("content-view css file: ", cssFileName)
		var ferr error
		css, ferr = os.Create(cssFileName)
		if ferr != nil {
			log.Fatal("unable to create content-view css file: ", cssFileName, ferr)
		}
		defer css.Close()
	}

	// create template input
	td := TemplateData{fm, data}

	//
	// create template, inject helper functions
	//
	templatePath := "./templates/"
	templateName := "content_view.gohtml"
	t, err := template.New(templateName).Funcs(template.FuncMap{
		"isLeaf": isLeaf, "isArray": isArray,
		"isObject": isObject, "isObjectArray": isObjectArray,
	}).ParseFiles(templatePath + templateName)
	if err != nil {
		return errors.Wrap(err, "error parsing content template ")
	}
	// execute the template
	err = t.Execute(html, td)
	if err != nil {
		return errors.Wrap(err, "error executing content template ")
	}

	//
	// now create the css file reflecting the
	// nesting structure of the json
	//
	if genCss {
		//
		// reduce the json to structural keys
		//
		flatd := Flatten(data)
		keys := make([]string, 0)
		for k, _ := range flatd {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		//
		// write css file
		//
		css.WriteString(contentCssBoilerplate())
		css.WriteString("\n\n/* auto-generated css from json structure below */\n\n")
		for _, s := range keys {
			str := s
			removeTool := strings.ReplaceAll(str, ":tool", "")
			splitPredicates := strings.ReplaceAll(removeTool, ".", " > .")
			cssLine := fmt.Sprintf(".%s { display: none; }\n\n", splitPredicates)
			css.WriteString(cssLine)
		}
		log.Println("css keys: ", len(flatd))
	}

	//
	// copy the supporting css & javascript files
	//
	var fileError error
	//
	// the json-oriented css file
	//
	srcFile := fmt.Sprintf("%s/css/content.css", resourcesFolder)
	dstFile := fmt.Sprintf("%s/%s/content.css", baseOutputFolder, fm.CssPath)
	fileError = copyFileContents(srcFile, dstFile)
	//
	// overrirdes file for above - keeps changes simpler
	//
	srcFile = fmt.Sprintf("%s/css/content_overrides.css", resourcesFolder)
	dstFile = fmt.Sprintf("%s/%s/content_overrides.css", baseOutputFolder, fm.CssPath)
	fileError = copyFileContents(srcFile, dstFile)

	//
	// audit file javascript - mostly toggling css attributes
	//
	srcFile = fmt.Sprintf("%s/js/content.js", resourcesFolder)
	dstFile = fmt.Sprintf("%s/%s/content.js", baseOutputFolder, fm.JsPath)
	fileError = copyFileContents(srcFile, dstFile)

	return fileError

}

//
// creates searchable glossary page for the syllabus
// Glossaries can exist at root (KLA) or syllabus level
// syllabus takes precedence, so is checked first, if not present
// kla-level is returned
//
func createGlossaryView(fm FileMeta, data interface{}) error {
	// create the output file
	outFileName := fmt.Sprintf("%s%s/%s", baseOutputFolder, fm.RootName, "glossary.html")
	log.Println("glossary-view output file: ", outFileName)
	html, err := os.Create(outFileName)
	if err != nil {
		log.Fatal("unable to create glossary-view html file: ", outFileName)
	}
	defer html.Close()

	// create template input
	td := TemplateData{fm, data}

	//
	// create template, inject helper functions
	//
	templatePath := "./templates/"
	templateName := "glossary_view.gohtml"
	t, err := template.New(templateName).Funcs(template.FuncMap{
		"isLeaf": isLeaf, "isArray": isArray,
		"isObject": isObject, "isObjectArray": isObjectArray,
	}).ParseFiles(templatePath + templateName)
	if err != nil {
		return errors.Wrap(err, "error parsing glossary template ")
	}
	// execute the template
	err = t.Execute(html, td)
	if err != nil {
		return errors.Wrap(err, "error executing glossary template ")
	}

	//
	// copy the supporting css & javascript files
	//
	var fileError error
	//
	// the json-oriented css file
	//
	srcFile := fmt.Sprintf("%s/css/glossary.css", resourcesFolder)
	dstFile := fmt.Sprintf("%s/%s/glossary.css", baseOutputFolder, fm.CssPath)
	fileError = copyFileContents(srcFile, dstFile)
	//
	// audit file javascript - mostly toggling css attributes
	//
	srcFile = fmt.Sprintf("%s/js/glossary.js", resourcesFolder)
	dstFile = fmt.Sprintf("%s/%s/glossary.js", baseOutputFolder, fm.JsPath)
	fileError = copyFileContents(srcFile, dstFile)

	return fileError

}

//
// creates the full audit view of the document -
// shows all elements with structure css toggles.
//
// fm 	-	FileMeta object to allow templating of all paths etc.
// data -	json data parsed as interface{}
//
func createAuditView(fm FileMeta, data interface{}) error {

	// create the output file
	outFileName := fmt.Sprintf("%s%s/%s", baseOutputFolder, fm.RootName, "audit.html")
	log.Println("audit-view output file: ", outFileName)
	html, err := os.Create(outFileName)
	if err != nil {
		log.Fatal("unable to create audit_view html file: ", outFileName)
	}
	defer html.Close()

	// create template input
	td := TemplateData{fm, data}

	//
	// create template, inject helper functions
	//
	templatePath := "./templates/"
	templateName := "audit_view.gohtml"
	t, err := template.New(templateName).Funcs(template.FuncMap{
		"isLeaf": isLeaf, "isArray": isArray,
		"isObject": isObject, "isObjectArray": isObjectArray,
	}).ParseFiles(templatePath + templateName)
	if err != nil {
		return errors.Wrap(err, "error parsing audit template ")
	}
	// execute the template
	err = t.Execute(html, td)
	if err != nil {
		return errors.Wrap(err, "error executing audit template ")
	}

	//
	// copy the supporting css & javascript files
	//
	var fileError error
	//
	// the json-oriented css file
	//
	srcFile := fmt.Sprintf("%s/css/audit.css", resourcesFolder)
	dstFile := fmt.Sprintf("%s/%s/audit.css", baseOutputFolder, fm.CssPath)
	fileError = copyFileContents(srcFile, dstFile)
	//
	// audit file javascript - mostly toggling css attributes
	//
	srcFile = fmt.Sprintf("%s/js/audit.js", resourcesFolder)
	dstFile = fmt.Sprintf("%s/%s/audit.js", baseOutputFolder, fm.JsPath)
	fileError = copyFileContents(srcFile, dstFile)

	return fileError
}

//
// looks for .json files in the input folder
//
func parseInputFolder(folderName string) ([]string, error) {

	files := make([]string, 0)

	jsonFiles, _ := filepath.Glob(folderName + "/*.json")

	files = append(files, jsonFiles...)
	if len(files) == 0 {
		return nil, errors.New("No input data *.json files found in input folder.")
	}

	return files, nil
}

//
// construct all required paths etc. from the
// input file name
//
func buildFileMeta(fileName string) FileMeta {

	fm := FileMeta{}

	fm.FileName = fileName
	shortName := filepath.Base(fileName)
	fm.RootName = strings.TrimSuffix(shortName, ".json")
	fm.HtmlPath = fmt.Sprintf("%s/", fm.RootName)
	fm.CssPath = fmt.Sprintf("%s/css/", fm.RootName)
	fm.JsPath = fmt.Sprintf("%s/js/", fm.RootName)

	return fm

}

//
// based on the metadata, make sure all
// necessary output folders are created
//
func makeOutputFolders(fm FileMeta) error {

	err := os.MkdirAll(baseOutputFolder+fm.HtmlPath, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.MkdirAll(baseOutputFolder+fm.CssPath, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.MkdirAll(baseOutputFolder+fm.JsPath, os.ModePerm)
	if err != nil {
		return err
	}

	return nil

}

//
// file copy
//
func copyFileContents(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	return err
}

//
// Copy directory from specified source path to destination path
//
func copyDirectory(scrDir, dest string) error {
	entries, err := ioutil.ReadDir(scrDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(scrDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		err := copyFileContents(sourcePath, destPath)
		if err != nil {
			fmt.Println("Ignored invalid file ", sourcePath)
		}

	}
	return nil
}

//
// Flatten takes a map of a json file and returns a new one where nested maps are replaced
// by dot-delimited keys.
//
func Flatten(m interface{}) map[string]interface{} {
	o := make(map[string]interface{})
	for k, v := range m.(map[string]interface{}) {
		switch child := v.(type) {
		case map[string]interface{}: // nested map (object)
			nm := Flatten(child)
			for nk, nv := range nm {
				o[k+"."+nk] = nv
			}
		case []interface{}: // array
			for _, sv := range child {
				switch member := sv.(type) {
				case map[string]interface{}: // object inside array
					nm := Flatten(member)
					for nk, nv := range nm {
						key := fmt.Sprintf("%s.%s", k, nk)
						o[key] = nv
					}
				default:
					key := fmt.Sprintf("%s", k)
					o[key] = sv
				}
			}
		default:
			o[k] = v
		}
	}
	return o
}

//
// simple css helpers to always add to the
// content css before the generated structural
// elements
//
func contentCssBoilerplate() string {
	return `

* {
    margin: 0;
    padding: 0;
}

body {
    font-family: 'PT Sans Narrow', '';
    font-size: 16px;
    color: #555;
}

a,
a:visited,
a:active {
    color: #1d77c2;
    text-decoration: none;
}

a:hover {
    color: #555;
}

main {
    margin: 3em auto 6em auto;
    max-width: 850px;
    line-height: 1.6em;
}

h1 {
    font-size: 3em;
    font-weight: normal;
    margin: 0 0 1em;
}

h2 {
    color: #1d77c2;
    font-size: 1.8em;
    font-weight: normal;
    margin: 3em 0 1em;
}

h3 {
    font-size: 1.2em;
    font-weight: bold;
    margin: 2em 0 1em;
}

h4 {
    font-size: 1.2em;
    font-weight: bold;
    margin: 2em 0 1em;
}

li {
    list-style-type: none;
}

/* minimap positioning */
#map {
    position: fixed;
    top: 0;
    right: 0;
    width: 160px;
    height: 100%;
    z-index: 100;
}



`
}
