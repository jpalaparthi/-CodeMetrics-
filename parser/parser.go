//Node parser
//Author: JP
//Purpose: To read all JS files and segrigate required comments and functions into JSON format.

// package parser is used to parse the file and give output as mentioned
// The current parser has been made only to read functions, number of lines of code
// And also for Number of chars used , all comments realted stuff.
// Yet to look for all files in the folder as this is only written for one file.
// Yet to convert it to the JSON format.
package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type Parser interface {
	Parse() []string
}

type parser struct {
	FileName     string   `json:fileName`
	LineCount    int      `json:lineCount`
	CharCount    int      `json:charsCount`
	FileData     []string `json:"-"`
	Comments     []string `json:comments`
	Functions    []string `json:functions`
	Dependencies []string `json:dependencies`
}

func New(filename string) (p *parser, err error) {
	p = &parser{FileName: filename}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	p.FileData = strings.Split(string(b), "\n")
	p.LineCount = len(p.FileData)
	return p, nil
}

func (p *parser) GetCharCount() int {
	p.CharCount = 0
	for _, v := range p.FileData {
		p.CharCount = p.CharCount + len(v)
	}
	return p.CharCount
}

func (p *parser) GetFuncs(line string) string {
	if strings.Contains(line, "function") && (!strings.Contains(line, "//") || !strings.Contains(line, "/*") || !strings.Contains(line, "*")) {
		line = strings.Replace(line, "{", "", 1)
		return line
	} else {
		return ""
	}
}

func (p *parser) GetDependencies(line string) string {
	line = strings.Replace(line, " ", "", -1) //strings.TrimSpace(line)
	line = strings.Replace(line, "'", "\"", -1)
	if strings.Contains(line, "=require(") && (!strings.Contains(line, "//") || !strings.Contains(line, "/*") || !strings.Contains(line, "*")) {
		n := strings.Index(line, "=require(\"")
		n = n + 10
		n1 := strings.LastIndex(line, "\")")
		if n1 > 0 {
			line = line[n:n1]
		} else {
			return ""
		}
		return line
	} else {
		return ""
	}
}

func (p *parser) Parse() []string {
	comments := make([]string, 0)
	cc := 0
	cc1 := 0
	cOn := false
	for _, v := range p.FileData {
		// To add functions
		if p.GetFuncs(v) != "" {
			p.Functions = append(p.Functions, p.GetFuncs(v))
		}

		// TO add dependencies
		if p.GetDependencies(v) != "" {
			p.Dependencies = append(p.Dependencies, p.GetDependencies(v))
		}
		// Chars count
		p.CharCount = p.CharCount + len(v)

		//Another comments
		if strings.Contains(v, "/*") {
			cOn = true
			cc1 = 0
		}
		if cOn {
			comments = append(comments, v)
			if strings.Contains(v, "*/") {
				cOn = false
				cc1++
			}
		} else if cc1 >= 1 {
			if v == "" {
				comments = append(comments, v)
			} else if v != "" {
				fmt.Println("Yes 2" + v)
				//str := strings.Replace(v, "{", "", -1)
				comments = append(comments, v)
				cc1 = 0
			}

		}

		// Comments part
		if strings.Contains(v, "//") && !cOn {
			comments = append(comments, v)
			cc++
		} else {
			if cc >= 1 {
				if v == "" {
					comments = append(comments, v)
				} else if v != "" {
					fmt.Println("Yes 2" + v)
					//str := strings.Replace(v, "{", "", -1)
					comments = append(comments, v)
					cc = 0
				}

			}
		}
	}
	return comments
}

func (p *parser) GenerateJSON() (string, error) {
	p.Comments = p.Parse()
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
