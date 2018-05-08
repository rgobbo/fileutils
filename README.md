# fileutils
Go (golang) library with helper functions to manipulate files.

## Installation

As a library

```shell
go get github.com/rgobbo/fileutils
```

## Usage

In your Go app you can do something like

```go
package main

import (
    "github.com/rgobbo/fileutils"
    "log"
    "os"
)

type Config struct {
    	HTTP          string         `json:"http.port"`
    	HTTPS         string         `json:"https.port"`
    	Host          string         `json:"host"`
    	TemplateDir   string         `json:"template.dir"`
    	StaticDir     string         `json:"http.static"`
    	HTTPError     string         `json:"http.error"`
    	DBType        string         `json:"db.type"`
    	DBServer      string         `json:"db.server"`
    	DBName        string         `json:"db.name"`
    	DBuser        string         `json:"db.user"`
    	DBPass        string         `json:"db.pass"`
}

func main() {
  confPath := "./config.json"
  var conf Config

  err := fileutils.LoadJson(confPath, &conf)
  if err != nil {
  	log.Fatal("Error loading json :", err)
  }

  log.Println("File processed successfully !!"


}
```

Or you can use a map[string] interface{} to load json files:

```go
    var conf map[string]interface{}

      err := fileutils.LoadJson(confPath, &conf)
      if err != nil {
      	log.Fatal("Error loading json :", err)
      }
```


## Documentation

###LoadByteFiles
 - Read files from folder and return bytes, filtered by extension

###LoadListFiles
 - Loads a directory recursivily and returns a list
 path : directory path to start create the list
 ext : extension to filter, only files with a especific extension will be included into list , other files will bo ignored
 removeExtension : if true remove extension form the file name
 Example : LoadListFiles("/Users/test", ".html", true)

###LoadJson
 - Load a json file and return into a inrterface
 Example : LoadJson("./file.json", &obj)

###SaveJson
 - Convert a interface into json and save a file

###LoadYaml
 - Load yaml file into a interface{}
 Example : LoadYaml("./file.yaml", &obj)

###SaveYaml
 - Save yaml interface{} into a file

###Unzip
 - Unizp a file to defined path
 Example: Unzip("/tmp/report-2015.zip", "/tmp/reports/")

###Zipit
 - Zip a directory or a file into a zip file
 Example : Zipit("/tmp/documents", "/tmp/backup.zip")
 Example : Zipit("/tmp/report.txt", "/tmp/report-2015.zip")

###GetCWD
 - Return working dir

###RenameIfExists
 - Rename a file if exists

###GetTimeStamp
 - Return timeStamp string with current date
