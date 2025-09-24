package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"unicode"

	"github.com/vorlif/spreak/v2/catalog/cldrplural"
)

type dataSet struct {
	Languages []string
	Rules     []*ruleData
}

type ruleData struct {
	Raw      string
	Category string
}

func (d *dataSet) Categories() []string {
	categories := make([]string, len(d.Rules))
	for i, r := range d.Rules {
		categories[i] = r.Category
	}
	return categories
}

func (d *dataSet) Name() string {
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

func (rd *ruleData) WithoutExamples() string {
	rawRule := rd.Raw[:strings.Index(rd.Raw, "@")]
	return strings.TrimSpace(rawRule)
}

type ruleJSON struct {
	Zero  string `json:"pluralRule-count-zero"`
	One   string `json:"pluralRule-count-one"`
	Two   string `json:"pluralRule-count-two"`
	Few   string `json:"pluralRule-count-few"`
	Many  string `json:"pluralRule-count-many"`
	Other string `json:"pluralRule-count-other"`
}

func (j *ruleJSON) ToData() []*ruleData {
	data := make([]*ruleData, 0, 3)
	if j.Zero != "" {
		data = append(data, &ruleData{Raw: j.Zero, Category: cldrplural.CategoryNames[cldrplural.Zero]})
	}
	if j.One != "" {
		data = append(data, &ruleData{Raw: j.One, Category: cldrplural.CategoryNames[cldrplural.One]})
	}
	if j.Two != "" {
		data = append(data, &ruleData{Raw: j.Two, Category: cldrplural.CategoryNames[cldrplural.Two]})
	}
	if j.Few != "" {
		data = append(data, &ruleData{Raw: j.Few, Category: cldrplural.CategoryNames[cldrplural.Few]})
	}
	if j.Many != "" {
		data = append(data, &ruleData{Raw: j.Many, Category: cldrplural.CategoryNames[cldrplural.Many]})
	}
	if j.Other != "" {
		data = append(data, &ruleData{Raw: j.Other, Category: cldrplural.CategoryNames[cldrplural.Other]})
	}
	return data
}

func (j *ruleJSON) hash() string {
	h := sha256.New()
	h.Write([]byte(j.Zero))
	h.Write([]byte(j.One))
	h.Write([]byte(j.Two))
	h.Write([]byte(j.Few))
	h.Write([]byte(j.Many))
	h.Write([]byte(j.Other))
	return hex.EncodeToString(h.Sum(nil))
}

type pluralsFile struct {
	Supplemental struct {
		Version struct {
			UnicodeVersion string `json:"_unicodeVersion"`
			CldrVersion    string `json:"_cldrVersion"`
		} `json:"version"`
		PluralsTypeCardinal map[string]ruleJSON `json:"plurals-type-cardinal"`
	} `json:"supplemental"`
}
