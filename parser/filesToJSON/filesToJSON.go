package filesToJSON

import (
	"encoding/json"
	"errors"
	_ "fmt"
	"io/ioutil"
	"os"
	"parser"
	"path/filepath"
	"strings"
)

type filesToJSON struct {
	FileList       []string
	ConfigFields   map[string][]string
	TotLineCount   int
	TotCharCount   int
	TotNoOfFiles   int
	Projectname    string
	Projectversion string
	FullJSON       string
}

func New(dir, projectName, projectVersion string) (a *filesToJSON, err error) {
	if projectName == "" || projectVersion == "" {
		return nil, errors.New("project name or project version is not given")
	}
	a = &filesToJSON{Projectname: projectName, Projectversion: projectVersion}

	a.ConfigFields, err = a.ReadConfig("config.json")
	if err != nil {
		return nil, err
	}
	exts := a.ConfigFields["accept-ext"]
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			//fmt.Println(f.Name())
		}
		if (a.SearchString(exts, filepath.Ext(path))) && (!f.IsDir()) {
			a.FileList = append(a.FileList, path)
			a.TotNoOfFiles++
			//fmt.Println(path)
		}
		return nil
	})
	a.Projectname = projectName
	a.Projectversion = projectVersion

	return a, nil
}

func (a *filesToJSON) ToJSONArray(exts []string) (string, error) {
	json := "\"files\":["
	for _, file := range a.FileList {
		p, err := parser.New(file)
		if err != nil {
			return "", err
		}
		fJSON, err := p.GenerateJSON()
		a.TotLineCount = a.TotLineCount + p.LineCount
		a.TotCharCount = a.TotCharCount + p.CharCount
		if err != nil {
			return "", err
		} else {
			json = json + fJSON + ","
		}
	}
	return json[0:len(json)-1] + "]", nil
}

func (a *filesToJSON) ToFullJSON(JSONStrings ...string) string {
	var JSON string

	for _, v := range JSONStrings {
		JSON = JSON + string(v)
	}
	return JSON
}

func (a *filesToJSON) WriteTOFile(filename, data string) (n int, err error) {
	f, err := os.Create(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	n, err = f.WriteString(data)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (a *filesToJSON) ReadConfig(filename string) (map[string][]string, error) {
	buf, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	fields := map[string][]string{}
	dec := json.NewDecoder(strings.NewReader(string(buf)))
	err = dec.Decode(&fields)
	if err != nil {
		return nil, errors.New("error in config.json file," + err.Error())
	}
	return fields, nil
}

func (a *filesToJSON) SearchString(list []string, str string) bool {
	for _, v := range list {
		if strings.ToUpper(v) == strings.ToUpper(str) {
			return true
		}
	}
	return false
}
