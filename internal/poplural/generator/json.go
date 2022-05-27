package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/vorlif/spreak/internal/poplural/ast"
)

var pluralsFile = filepath.FromSlash(filepath.Join("./plurals.json"))

func getJSONRules() []*RuleData {
	data, err := os.ReadFile(pluralsFile)
	checkError(err)

	root := PluralsFile{}
	checkError(json.Unmarshal(data, &root))

	rulesData := map[string]*RuleData{}
	for lang, fileData := range root {
		if rd, ok := rulesData[fileData.Formula]; ok {
			rd.Languages = append(rd.Languages, lang)
			continue
		}

		formula := fmt.Sprintf("nplurals=%d; plural=%s;", fileData.Plurals, fileData.Formula)
		parsed := ast.MustParse(formula)
		rulesData[fileData.Formula] = &RuleData{
			Languages:   []string{lang},
			Raw:         formula,
			Count:       fileData.Plurals,
			Compiled:    compileForms(parsed),
			CompiledRaw: ast.CompileToString(parsed),
		}
	}

	res := make([]*RuleData, 0, len(rulesData))
	for _, rd := range rulesData {
		res = append(res, rd)
	}

	return res
}

type PluralsFile map[string]FileLangData

type FileLangData struct {
	Language string
	Formula  string
	Plurals  int
	Cases    []string
	Examples map[string]string
}

type RuleData struct {
	Raw         string
	CompiledRaw string
	Compiled    string
	Count       int
	Languages   []string
}

func (d *RuleData) Name() string {
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
