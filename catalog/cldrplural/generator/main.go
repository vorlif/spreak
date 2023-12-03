package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/catalog/cldrplural/ast"
	"github.com/vorlif/spreak/internal/util"
)

var pluralsFilePath = filepath.FromSlash(filepath.Join("./plurals.json"))

func main() {
	data, err := os.ReadFile(pluralsFilePath)
	checkError(err)

	root := &pluralsFile{}
	checkError(json.Unmarshal(data, root))

	duplicateMap := make(map[string]*dataSet, 40)
	dataSets := make([]*dataSet, 0, 40)
	for lang, jsonRules := range root.Supplemental.PluralsTypeCardinal {
		h := jsonRules.hash()
		if _, ok := duplicateMap[h]; ok {
			duplicateMap[h].Languages = append(duplicateMap[h].Languages, lang)
			continue
		}

		ds := &dataSet{
			Languages: []string{lang},
			Rules:     jsonRules.ToData(),
		}
		duplicateMap[h] = ds
		dataSets = append(dataSets, ds)
	}

	for _, dataset := range dataSets {
		sort.Strings(dataset.Languages)
	}

	sort.Slice(dataSets, func(i, j int) bool {
		if a, b := len(dataSets[i].Rules), len(dataSets[j].Rules); a != b {
			return a < b
		}

		return dataSets[i].Name() < dataSets[j].Name()
	})

	executeAndSafe("../builtin_gen.go", codeTemplate, dataSets)
	executeAndSafe("../builtin_gen_test.go", testTemplate, dataSets)
	executeAndSafe("../evaluation_gen_test.go", evaluationTest, dataSets)
}

// executeAndSafe applies the dataSet's to a parsed template and saves the result correctly
// formatted in a file 'name'.
func executeAndSafe(name string, tmpl *template.Template, dataSets []*dataSet) {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, dataSets)
	checkError(err)

	p, err := format.Source(buf.Bytes())
	checkError(err)

	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	checkError(err)
	defer file.Close()

	_, err = file.Write(p)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
}

var templateFuncs template.FuncMap = map[string]interface{}{
	"ConvertToGo": func(rawRule string) string {
		rule, err := ast.Parse(rawRule)
		checkError(err)
		return compileNode(rule.Root)
	},
	"ExtractSamples": func(rawRule string) []string {
		rule, err := ast.Parse(rawRule)
		checkError(err)
		return rule.Samples
	},
	"HasCondition": func(rawRule string) bool {
		return !strings.HasPrefix(strings.TrimSpace(rawRule), "@")
	},
	"CaseExpression": func(langs []string) string {
		result := make(map[string]struct{}, len(langs))

		for _, lang := range langs {
			result[fmt.Sprintf(`"%s"`, lang)] = struct{}{}

			// For some strings, the parsed tag is not identical to the string.
			// These cases are checked separately here.
			parsedLang := language.MustParse(lang).String()
			if lang != parsedLang {
				result[fmt.Sprintf(`"%s"`, parsedLang)] = struct{}{}
			}
		}

		languageStrings := util.Keys(result)
		sort.Strings(languageStrings)
		return strings.Join(languageStrings, ",")
	},
}

var codeTemplate = template.Must(template.New("rule").Funcs(templateFuncs).Parse(`// This file is generated by cldrplural/generator/generate.sh; DO NOT EDIT
package cldrplural

import (
	"math"
	"slices"
)

func forLanguage(lang string) *RuleSet {
	switch lang {
	{{- range $key, $data := .}}
	case {{CaseExpression $data.Languages}}:
		return newRuleSet{{$data.Name}}()
	{{- end}}
	default:
		return nil
	}
}

{{range $key, $data := .}}
// {{printf "%v" $data.Languages}}
func newRuleSet{{$data.Name}}() *RuleSet {
	return &RuleSet{
		Categories: []Category{{"{"}}{{range $i, $cat := $data.Categories}}{{if $i}}, {{end}}{{$cat}}{{end}}{{"}"}},
		FormFunc: func(ops *Operands) Category { 
			{{- range $data.Rules}}
				{{- if HasCondition .Raw}}
					// {{.WithoutExamples}}
					if {{ConvertToGo .Raw}} {
						return {{.Category}}
					}
				{{else}}
					{{if .WithoutExamples}}// {{.WithoutExamples}}{{end}}
					return {{.Category}}
				{{- end -}}
			{{end}}
		},
	}
}
{{end}}
`))

var testTemplate = template.Must(template.New("rule").Funcs(templateFuncs).Parse(`// This file is generated by cldrplural/generator/generate.sh; DO NOT EDIT
package cldrplural

import (
	"testing"

	"golang.org/x/text/language"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

{{range $key, $data := .}}
func TestBuiltin{{$data.Name}}(t *testing.T) {
	for _, lang := range {{printf "%#v" $data.Languages}} {
		set := forLanguage(language.MustParse(lang).String())
		require.NotNilf(t, set, "RuleSet for language %s (%s) not found", lang, language.MustParse(lang).String())
		{{range $data.Rules}}
			{{$samples := ExtractSamples .Raw}}
			for _, sample := range {{printf "%#v" $samples}} {
				op := MustNewOperands(sample)
				assert.Equal(t, {{.Category}}, set.FormFunc(op)) 
			}
		{{end}}
	}
}
{{end}}
`))

var evaluationTest = template.Must(template.New("rule").Funcs(templateFuncs).Parse(`// This file is generated by cldrplural/generator/generate.sh; DO NOT EDIT
package cldrplural

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vorlif/spreak/catalog/cldrplural/ast"
)

{{range $key, $data := .}}
func TestEvaluate{{$data.Name}}(t *testing.T) {
	var rule *ast.Rule
	{{range $data.Rules}}
		rule = ast.MustParse("{{.Raw}}")
		
		{{$samples := ExtractSamples .Raw}}
		for _, sample := range {{printf "%#v" $samples}} {
			op := MustNewOperands(sample)
			assert.True(t, evaluate(rule, op), sample)
		}
	{{end}}
}
{{end}}
`))
