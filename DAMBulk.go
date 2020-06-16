package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	//   	"time"
	/* 	"github.com/Tkanos/gonfig"              // config management support
	   	"github.com/lib/pq"                     // golang postgres db driver
	*/)

// wrapper around linux grep utility
func logMessage(message string, templatename string, errortype string) {
	//log.Fatal(message + " " + templatename + " " + errortype)
	println(message + " " + templatename + " " + errortype)
}

func findFilename(pattern string, dir string) string {
	cmd := exec.Command("find", dir, "-name ", pattern)

	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Run()

	stdout := outbuf.String()
	return stdout
}

func grepDir(pattern string, dir string) string {
	cmd := exec.Command("grep", "-r", "--exclude-dir=downloads", pattern, dir)

	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Run()

	stdout := outbuf.String()
	return stdout
}
func getTemplateID(filepath string) string {

	templateID := ""
	file, err := os.Open(filepath)
	if err != nil {
		logMessage("getTemplateID os.Open() failed : "+err.Error(), "", "ERROR")
		return templateID
	}

	lines, err := readFile(file)
	if err != nil {
		logMessage("getTemplateID readFile() failed : "+err.Error(), "", "ERROR")
		return templateID
	}
	content := strings.Join(lines, " ")
	splits := strings.Split(strings.ToLower(content), "<id>")
	if len(splits) > 1 {
		templateID = (splits[1])[0:36] // TODO: HACK: assumes id format is fixed....
	}
	defer file.Close()

	return templateID

}

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
func readlines2(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	r := bufio.NewReader(file)
	s, e := Readln(r)
	for e == nil {

		lines = append(lines, s)
		s, e = Readln(r)
	}

	return lines, nil
}
func readFile(theFile *os.File) ([]string, error) {

	var lines []string
	r := bufio.NewReader(theFile)
	s, e := Readln(r)
	for e == nil {

		lines = append(lines, s)
		s, e = Readln(r)
	}

	return lines, nil
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func getContents(inputString string, delimeter1 string, delimeter2 string) string {
	var result = ""
	var a = strings.Index(inputString, delimeter1)
	var b = strings.Index(inputString, delimeter2)
	if a != -1 && b != -1 && a < b {
		result = inputString[a+len(delimeter1) : b]
	}
	return result
}

func main() {
	fmt.Println("DAMBulk2")
	//exePath := os.Args[0]

	argsWithProg := os.Args
	fmt.Println(argsWithProg)

	argsWithoutProg := os.Args[1:]
	fmt.Println(argsWithoutProg)

	/* 	configPath, exeName := filepath.Split(exePath)

	   	fmt.Println( "exename = " + exeName )
	*/ //	fmt.Println( "configPath = " + configPath + configFilename)

	// for every .oet in the directory,
	// 	- extract its id
	//  - search mirror for id
	//  - if found, update the .oet with the ckm mirror id and write it out to a new directory

	gendir := "./generated"
	nomatchdir := "./nomatch"
	matchdir := "./match"

	files, err := ioutil.ReadDir(gendir)
	if err != nil {
		log.Fatal(err)
	}

	currentdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	templateIdMap := make(map[string]string)

	//We need to create a hashmap of names with id's
	//Loop through all the files in  /opt/ckm-mirror/local/templates/entry/instruction/
	//For each file extract the  <name></name><id></id>
	//mirrorFilePath := "/opt/ckm-mirror/local/templates/entry/instruction/"
	/*
		mirrorfiles, err := ioutil.ReadDir(mirrorFilePath)
		if err != nil {
			log.Fatal(err)
		}

		for _, mirrorfile := range mirrorfiles {
			fmt.Println("Examining mirrorfile " + mirrorfile.Name())
			var fullpath = mirrorFilePath + mirrorfile.Name()
			content, err := ioutil.ReadFile(fullpath)
			if err != nil {
				log.Fatal(err)
			}
			stringcontents := string(content)
			var templateName = getContents(stringcontents, "<name>", "</name>")
			var templateId = getContents(stringcontents, "<id>", "</id>")
			fmt.Println("Found templateName= " + templateName + " and id=" + templateId)
			templateIdMap[templateName] = templateId
		}
	*/

	//Get files recursively from some directory
	var numberFilesRead = 0
	mirrorFilePath := "/opt/ckm-mirror/local/templates/"
	err2 := filepath.Walk(mirrorFilePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			//fmt.Println(path, info.Size())
			if fi, err := os.Stat(path); err == nil {
				if fi.Mode().IsDir() {
					//fmt.Println("Is a Directory")
				} else if strings.HasSuffix(strings.ToUpper(""+path), ".OET") {
					//fmt.Println("filename=" + fi.Name())
					//var fullpath = template_path + "\\" + path
					numberFilesRead++
					fmt.Println("Examining " + path)
					//var outputFile = mirrorFilePath
					// + fi.Name()
					//fmt.Println("Output Path = " + outputFile)
					//removeSurplusArchetypes(path, outputFile)
					content, err := ioutil.ReadFile(path)
					if err != nil {
						log.Fatal(err)
					}
					stringcontents := string(content)
					var templateName = getContents(stringcontents, "<name>", "</name>")
					var templateId = getContents(stringcontents, "<id>", "</id>")
					fmt.Println("Found templateName= " + templateName + " and id=" + templateId)
					templateIdMap[templateName] = templateId
				}
			}

			return nil
		})
	if err2 != nil {
		log.Println(err)
	}

	fmt.Println("# items in templateIdMap = ")
	fmt.Println(len(templateIdMap))

	//Display it afterwards
	fmt.Println(currentdir) //Reading the contents of the file
	fmt.Println("Going to examine the files")
	for _, file := range files {

		//DAMBulk is using the file name to match to the mirror
		//It should be extracting the <name> value from the file and using that to look for the mirror
		//Get the current executable folder

		fmt.Println(file.Name())
		if filepath.Ext(file.Name()) == ".oet" {

			fmt.Println("Examining " + file.Name())
			var fullpath = currentdir + "/generated/" + file.Name()
			fmt.Println("FullPath= " + fullpath)
			content, err := ioutil.ReadFile(fullpath)
			if err != nil {
				log.Fatal(err)
			}
			stringcontents := string(content)
			var templateName = getContents(stringcontents, "<name>", "</name>")
			fmt.Println("templateName=" + templateName)
			//			templateid := getTemplateID(gendir + "/" + file.Name())

			//			ckmtemplateid := getTemplateID("/opt/ckm-mirror/local/templates/entry/instruction/" + file.Name())
			//ckmtemplateid := getTemplateID("/opt/ckm-mirror/local/templates/entry/instruction/" + templateName + ".oet")
			idMatch := templateIdMap[templateName]
			fmt.Println("idMatch=" + idMatch)

			if idMatch == "" {
				// couldn't find in mirror
				// copy to not matched dir
				fmt.Println("was not a match")
				err := Copy(gendir+"/"+file.Name(), nomatchdir+"/"+file.Name())
				if err != nil {
					logMessage(err.Error(), file.Name(), "ERROR")
				}

			} else {

				fmt.Println("was a match")
				newfile := matchdir + "/" + file.Name()
				err := Copy(gendir+"/"+file.Name(), newfile)
				if err != nil {
					logMessage(err.Error(), file.Name(), "ERROR")
				}

				// replace the newfile template id with the one from the mirror
				read, err := ioutil.ReadFile(newfile)
				if err != nil {
					panic(err)
				}

				fmt.Println(newfile + " and updating id to be " + idMatch)
				gentemplateID := getTemplateID(newfile)

				newContents := strings.Replace(string(read), gentemplateID, idMatch, -1)

				//fmt.Println(newContents)

				err = ioutil.WriteFile(newfile, []byte(newContents), 0)
				if err != nil {
					panic(err)
				}

			}

			//fmt.Println(ckmtemplateid)

			/* 			matches, err := filepath.Glob( "/opt/ckm-mirrir/local/templates/entry/instruction/" + file.Name() )
			   			if (err == nil) {

			   			}
			*/
			//			found := findFilename( file.Name(), "/opt/ckm-mirror/local/templates")
			//found := grepDir( templateid, "/opt/ckm-mirror/local/templates")
			//			fmt.Println( found )
		}
	}

	fmt.Println("# items in templateIdMap = ")
	fmt.Println(len(templateIdMap))

}

func findFile(targetDir string, pattern []string) {

	for _, v := range pattern {
		matches, err := filepath.Glob(targetDir + v)

		if err != nil {
			fmt.Println(err)
		}

		if len(matches) != 0 {
			fmt.Println("Found : ", matches)
		}
	}
}
