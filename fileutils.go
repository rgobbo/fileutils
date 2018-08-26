package fileutils

import (
	"encoding/json"
	"fmt"
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

//LoadFilesInfo - loads a directory recursivily and returns a map[string]interface - list of infos
// path : directory path to start create the list
// Example : LoadFilesInfo("/Users/test")
func LoadFilesInfo(path string) ([]map[string]interface{}, error) {

	list := []map[string]interface{}{}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return list, err
	}

	for _, f := range files {
		var listItem = make(map[string]interface{})
		if f.IsDir() {
			listItem["name"] = f.Name()
			abs := path + string(filepath.Separator) + f.Name()
			listItem["absolutePath"], _ = filepath.Abs(abs)
			listItem["extension"] = ""
			listItem["path"] = abs
			listItem["isDir"] = true
			listDir, err := LoadFilesInfo(path + "/" + f.Name())
			if err != nil {
				return list, err
			}
			listItem["childs"] = listDir

		} else {
			fileExt := filepath.Ext(f.Name())
			listItem["name"] = f.Name()
			abs := path + string(filepath.Separator) + f.Name()
			listItem["absolutePath"] = abs
			listItem["extension"] = fileExt[1:]
			listItem["path"] = abs
			listItem["isDir"] = false
			listItem["childs"] = nil
		}
		list = append(list, listItem)
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
	err = yaml.Unmarshal(file, obj)
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
