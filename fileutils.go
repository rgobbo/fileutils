package fileutils

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

//LoadByteFiles - read files from folder and return bytes, filtered by extension
func LoadByteFiles(dirname string, ext string) ([]byte, error) {
	var strCode []byte

	err := filepath.Walk(dirname, func(path string, f os.FileInfo, err error) error {
		// ext = .js
		if !f.IsDir() && strings.HasSuffix(f.Name(), ext) {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			strCode = append(strCode, b...)
		}
		return nil
	})
	return strCode, err
}

//LoadListFiles - loads a directory recursivily and returns a list
// path : directory path to start create the list
// ext : extension to filter, only files with a especific extension will be included into list , other files will bo ignored
// removeExtension : if true remove extension form the file name
// Example : LoadListFiles("/Users/test", ".html", true)
func LoadListFiles(path string, ext string, removeExtension bool) ([]string, error) {
	var list []string

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return list, err
	}

	for _, f := range files {
		if f.IsDir() {
			listDir, err := LoadListFiles(path+"/"+f.Name(), ext, removeExtension)
			if err != nil {
				return list, err
			}

			for _, s := range listDir {

				if removeExtension == true {
					list = append(list, f.Name()+"/"+strings.Replace(s, ext, "", 1))
				} else {
					list = append(list, f.Name()+"/"+s)
				}

			}

		} else {
			fileExt := filepath.Ext(f.Name())
			if fileExt == ext {
				if removeExtension == true {
					list = append(list, strings.Replace(f.Name(), ext, "", 1))
				} else {
					list = append(list, f.Name())
				}
			}
		}

	}

	return list, err
}

//LoadBytesDir - read files from folder
func LoadBytesDir(dirname string) ([]byte, error) {
	var strCode []byte

	err := filepath.Walk(dirname, func(path string, f os.FileInfo, err error) error {
		// ext = .js
		if !f.IsDir() {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			strCode = append(strCode, b...)
		}
		return nil
	})
	return strCode, err
}

//LoadJson  - Load a json file and return into a inrterface
//Example : LoadJson("./file.json", &obj)
func LoadJson(path string, obj interface{}) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &obj)
	if err != nil {
		return err
	}
	return nil
}

//SaveJson  - Convert a interface into json and save a file
func SaveJson(path string, obj interface{}) error {

	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, b, 0644)
}

//LoadYaml - Load yaml file into a interface{}
//Example : LoadYaml("./file.yaml", &obj)
func LoadYaml(path string, obj interface{}) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(file, &obj)
	if err != nil {
		return err
	}
	return nil
}

//SaveYaml - save yaml interface{} into a file
func SaveYaml(path string, obj interface{}) error {

	b, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, b, 0644)
}

//RemoveDuplicates - remove duplicate strings from slice string
func RemoveDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

//Unzip - unizp a file to defined path
//Unzip("/tmp/report-2015.zip", "/tmp/reports/")
func Unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
		fileReader.Close()
		targetFile.Close()
	}

	return nil
}

//Zipit - zip a directory or a file into a zip file
//Zipit("/tmp/documents", "/tmp/backup.zip")
//Zipit("/tmp/report.txt", "/tmp/report-2015.zip")
func Zipit(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

//GetCWD - return working dir
func GetCWD() string {
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR:  Could not get working directory.\n")
		fmt.Fprintf(os.Stderr, "ERROR-MESSAGE:%v\n", err)
		os.Exit(4)
	}
	return currentWorkingDirectory
}

//RenameIfExists - rename a file if exists
func RenameIfExists(path string) {
	os.Rename(path, fmt.Sprintf("%s-Pre-%s", path, GetTimeStamp()))
}

const TIME_LAYOUT = "Jan-02-2006_15-04-05-MST"

//GetTimeStamp - return timeStamp string with current date
func GetTimeStamp() string {
	now := time.Now()
	return now.Format(TIME_LAYOUT)
}
