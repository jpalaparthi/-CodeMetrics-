package main

import (
	"bufio"
	"fmt"
	"os"
	"parser/filesToJSON"
	"strings"
	"time"
)

func main() {
	var pn, pv, fn string
	fmt.Println("Project Name>")
	pn, err := bufio.NewReader(os.Stdin).ReadString('\n')
	pn = strings.Replace(pn, "\r\n", "", -1)
	fmt.Println("Project Version>")
	fmt.Scanln(&pv)
	fn = GetOutPutFileName()
	fmt.Println(pn, pv, fn)
	f, err := filesToJSON.New("./", pn, pv)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("-->received project name,version and output file details")
		fmt.Println("-->reading config.json file")
		fileexts := f.ConfigFields["accept-ext"]
		fmt.Println("-->converting data in files to json array")
		json, err := f.ToJSONArray(fileexts)
		fmt.Println("-->parsing files and writing to JSON array")
		if err != nil {
			fmt.Println(err)
		} else {
			jsonStart1 := "{\"projectName\":\"" + f.Projectname + "\",\"projectVersion\":\"" + f.Projectversion + "\","
			jsonStart2 := "\"totalNoOfLines\":" + fmt.Sprintf("%d", f.TotLineCount) + ",\"totalNoOfChars\":" + fmt.Sprintf("%d", f.TotCharCount) + ",\"totalNoOfFiles\":" + fmt.Sprintf("%d", f.TotNoOfFiles) + ",\"fileTakenOn\":\"" + time.Now().String() + "\","
			jsonEnd1 := "}"
			fullJSON := f.ToFullJSON(jsonStart1, jsonStart2, json, jsonEnd1)
			fmt.Println("-->writing output to ", fn, " file")
			f.WriteTOFile(fn, fullJSON)
			fmt.Println("-->\"" + fn + "\" output file has been created")
		}
	}
}

// Validate  output filer
func GetOutPutFileName() (fn string) {
	fmt.Println("Output File Name>")
	fmt.Scanln(&fn)
	if !strings.Contains(fn, ".json") {
		fmt.Println("Out put file should be .JSON file")
		return GetOutPutFileName()
	}
	return fn
}
