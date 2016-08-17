package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"parser"

	"path/filepath"
	"strings"
	"time"
)

func main1() {
	fileList := []string{}
	// To disply total count of lines of code from all modules.
	var totLineCount, totCharCount, totalNoOfFiles int
	//project name and project version
	var projectname, projectversion string
	projectname = "CodeMetrics"
	projectversion = "v0.1"
	jsonStart1 := "{\"projectName\":\"" + projectname + "\",\"projectVersion\":\"" + projectversion + "\","

	config, _ := ReadConfig("config.json")
	fileexts := config["accept-ext"]
	//acceptdirs:=config["accept-dir"]

	filepath.Walk("./", func(path string, f os.FileInfo, err error) error {
		//fmt.Println(path)
		fileList = append(fileList, path)
		return nil
	})

	var JSON string = jsonStart1 + "$jsonStart2$" + "\"files\":[" // Starting of json file
	for _, file := range fileList {
		fStat, _ := os.Stat(file)
		if (SearchString(fileexts, filepath.Ext(file))) && (!fStat.IsDir()) {
			totalNoOfFiles++
			p, _ := parser.New(file)
			fJSON, err := p.GenerateJSON()
			totLineCount = totLineCount + p.LineCount
			totCharCount = totCharCount + p.CharCount
			if err != nil {
				fmt.Println(err)
			} else {
				JSON = JSON + fJSON + ","
			}
		}
	}
	jsonStart2 := "\"totalNoOfLines\":" + fmt.Sprintf("%d", totLineCount) + ",\"totalNoOfChars\":" + fmt.Sprintf("%d", totCharCount) + ",\"totalNoOfFiles\":" + fmt.Sprintf("%d", totalNoOfFiles) + ",\"fileTakenOn\":\"" + time.Now().String() + "\","
	JSON = JSON[0:len(JSON)-1] + "]}"
	JSON = strings.Replace(JSON, "$jsonStart2$", jsonStart2, 1)
	WriteTOFile("output.json", JSON)

}

func WriteTOFile(filename, data string) {
	f, _ := os.OpenFile(filename, os.O_APPEND, 0667)
	f.WriteString(data)
	f.Close()
}

func ReadConfig(filename string) (map[string][]string, error) {
	buf, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	fields := map[string][]string{}
	dec := json.NewDecoder(strings.NewReader(string(buf)))
	dec.Decode(&fields)
	return fields, nil
}

func SearchString(list []string, str string) bool {
	for _, v := range list {
		if strings.ToUpper(v) == strings.ToUpper(str) {
			return true
		}
	}
	return false
}
