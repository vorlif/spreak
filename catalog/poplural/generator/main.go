package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"slices"
	"sort"
	"strings"
	"text/template"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/internal/util"
)

func main() {
	jsonRules := getJSONRules()
	extraRules := getExtraRules()

	rules := append(jsonRules, extraRules...)

	for _, r := range rules {
		sort.Strings(r.Languages)
	}

	sort.Slice(rules, func(i, j int) bool {
		if a, b := rules[i].Count, rules[j].Count; a != b {
			return a < b
		}

		return rules[i].Name() < rules[j].Name()
	})

	executeAndSafe("../builtin_gen.go", codeTemplate, rules)
	executeAndSafe("../builtin_gen_test.go", testTemplate, rules)
}

// executeAndSafe applies the DataSet's to a parsed template and saves the result correctly
// formatted in a file 'name'.
func executeAndSafe(name string, tmpl *template.Template, rules []*ruleData) {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, rules)
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

var templateFuncs template.FuncMap = map[string]any{
	"StrJoin": func(elems []string, sep string) string {
		values := slices.Clone(elems)
		for i := range values {
			values[i] = fmt.Sprintf(`"%s"`, values[i])
		}
		return strings.Join(values, ",")
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

var codeTemplate = template.Must(template.New("rule").Funcs(templateFuncs).Parse(`// This file is generated by poplural/generator/main.go; DO NOT EDIT
package poplural

func forRawRule(rawRule string) *Rule {
	switch rawRule {
	{{- range $key, $rule := .}}
	case "{{$rule.CompiledRaw}}":
		{{- if $rule.Languages}}
			return newRule{{$rule.Name}}()
		{{- else -}}
			return &Rule{
				NPlurals: {{$rule.Count}},
				// {{$rule.Raw}}
				FormFunc: func(n int64) int {
					{{$rule.Compiled}}
				},
			}
		{{- end -}}
	{{ end }}
	default:
		return nil
	}
}

func forLanguage(lang string) *Rule {
	switch lang {
	{{- range $key, $rule := .}}
	{{- if $rule.Languages -}}
	case {{CaseExpression $rule.Languages}}:
		return newRule{{$rule.Name}}()
	{{end -}}
	{{- end -}}
	default:
		return nil
	}
}

{{range $idx, $rule := .}}
	{{- if $rule.Languages -}}
		// {{printf "%v" $rule.Languages}}
		func newRule{{$rule.Name}}() *Rule {
			return &Rule{
				NPlurals: {{$rule.Count}},
				// {{$rule.Raw}}
				FormFunc: func(n int64) int {
					{{$rule.Compiled}}
				},
			}
		}
	{{end}}
{{end}}
`))

var testTemplate = template.Must(template.New("rule").Funcs(templateFuncs).Parse(`// This file is generated by cldrplural/generator/main.go; DO NOT EDIT
package poplural

import (
	"fmt"
	"testing"

	"golang.org/x/text/language"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

{{ range $key, $data := . }}
{{ if $data.Languages }}
func TestBuiltin{{$data.Name}}(t *testing.T) {
	for _, lang := range {{printf "%#v" $data.Languages}} {
		rule := forLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		{{range $result, $examples := $data.Examples }}
			for _, example := range {{printf "%#v" $examples}} {
				form, err := rule.Evaluate(example)
				require.NoError(t, err)
				assert.Equal(t, {{$result}}, form, fmt.Sprintf("rule.Evaluate(%s) should be {{$result}}", example))
			}
		{{end}}
	}
}
{{ end }}
{{ end }}
`))
