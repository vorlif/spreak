package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/vorlif/spreak/v2/catalog/poplural/ast"
)

var pluralsFilePath = filepath.FromSlash(filepath.Join("./plurals.json"))

func readJSONRules() []*ruleData {
	data, err := os.ReadFile(pluralsFilePath)
	checkError(err)

	root := pluralsFile{}
	checkError(json.Unmarshal(data, &root))

	rulesData := map[string]*ruleData{}
	for lang, fileData := range root {
		if rd, ok := rulesData[fileData.Formula]; ok {
			rd.Languages = append(rd.Languages, lang)
			continue
		}

		formula := fmt.Sprintf("nplurals=%d; plural=%s;", fileData.Plurals, fileData.Formula)
		parsed := ast.MustParse(formula)
		rulesData[fileData.Formula] = &ruleData{
			Languages:   []string{lang},
			Raw:         formula,
			Count:       fileData.Plurals,
			Compiled:    compileForms(parsed),
			CompiledRaw: ast.CompileToString(parsed),
			Examples:    extractExamples(fileData.Cases, fileData.Examples),
		}
	}

	res := make([]*ruleData, 0, len(rulesData))
	for _, rd := range rulesData {
		res = append(res, rd)
	}

	return res
}

type pluralsFile map[string]fileLangData

type fileLangData struct {
	Language string
	Formula  string
	Plurals  int
	Cases    []string
	Examples map[string]string
}

type ruleData struct {
	Raw         string
	CompiledRaw string
	Compiled    string
	Count       int
	Languages   []string
	// Maps from index to examples
	Examples map[int][]string
}

func (d *ruleData) Name() string {
	if len(d.Languages) == 0 {
		return d.CompiledRaw
	}

	var b strings.Builder
	for _, lang := range d.Languages {
		r := []rune(lang)
		r[0] = unicode.ToUpper(r[0])
		lang = string(r)
		lang = strings.ReplaceAll(lang, " ", "")
		lang = strings.ReplaceAll(lang, "-", "_")
		b.WriteString(lang)
	}

	return b.String()
}

func extractExamples(cases []string, data map[string]string) map[int][]string {

	extracted := make(map[int][]string, len(data))

	for caseName, caseExamples := range data {
		examples := make([]string, 0)

		for _, example := range strings.Split(caseExamples, ",") {
			// The last item is always an omission point
			if example == " …" {
				break
			}

			example = strings.TrimSpace(example)

			if parts := strings.Split(example, "~"); len(parts) == 2 {
				startNumber, err := strconv.Atoi(parts[0])
				checkError(err)
				stopNumber, err := strconv.Atoi(parts[1])
				checkError(err)

				for i := startNumber; i <= stopNumber; i++ {
					examples = append(examples, strconv.Itoa(i))
				}
			} else {
				// The JSON file uses "6c3" instead of "6e3".
				example = strings.ReplaceAll(example, "c", "e")
				examples = append(examples, example)
			}
		}

		// Since the data is stored in a map, the appropriate index must be searched for.
		for i, name := range cases {
			if name == caseName {
				extracted[i] = examples
				break
			}
		}
	}

	return extracted
}
