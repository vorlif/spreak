package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/vorlif/spreak/catalog/poplural/ast"
)

var pluralsFilePath = filepath.FromSlash(filepath.Join("./plurals.json"))

func getJSONRules() []*ruleData {
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
