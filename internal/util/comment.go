package util

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

type Reference struct {
	Path   string
	Line   int
	Column int
}

func (r Reference) String() string {
	if r.Line >= 0 {
		return fmt.Sprintf("%s:%d", r.Path, r.Line)
	}

	return r.Path
}

type Comment struct {
	Translator     string       // #  translator-comments
	Extracted      string       // #. extracted-comments
	References     []*Reference // #: src/file.go:338
	Flags          []string     // #, fuzzy,go-format,range:0..10
	PrevMsgContext string       // #| msgctxt previous-context
	PrevMsgID      string       // #| msgid previous-untranslated-string
}

func (p *Comment) AddReference(ref *Reference) {
	if p.References == nil {
		p.References = make([]*Reference, 0)
	}
	p.References = append(p.References, ref)
	p.sortReferences()
}

func (p *Comment) sortReferences() {
	sort.Slice(p.References, func(i, j int) bool {
		if p.References[i].Path == p.References[j].Path {
			return p.References[i].Line < p.References[j].Line
		}

		return p.References[i].Path < p.References[j].Path
	})
}

func (p *Comment) Less(q *Comment) bool {
	for i := 0; i < len(p.References); i++ {
		if i >= len(q.References) {
			break
		}
		if a, b := p.References[i].Path, q.References[i].Path; a != b {
			return a < b
		}
		if a, b := p.References[i].Line, q.References[i].Line; a != b {
			return a < b
		}
		if a, b := p.References[i].Column, q.References[i].Column; a != b {
			return a < b
		}
	}
	return false
}

func (p *Comment) HasFlag(flag string) bool {
	for _, s := range p.Flags {
		if s == flag {
			return true
		}
	}
	return false
}

func (p Comment) WriteTo(w io.StringWriter, wrapWidth int) error {
	if p.Translator != "" {
		for _, comment := range WrapString(p.Translator, wrapWidth) {
			if _, err := w.WriteString(fmt.Sprintf("# %s\n", comment)); err != nil {
				return err
			}
		}
	}

	if p.Extracted != "" {
		for _, comment := range WrapString(p.Extracted, wrapWidth) {
			if _, err := w.WriteString(fmt.Sprintf("#. %s\n", comment)); err != nil {
				return err
			}
		}
	}

	if len(p.References) > 0 {
		var buff bytes.Buffer
		for _, ref := range p.References {
			buff.WriteString(fmt.Sprintf("%s:%d ", ref.Path, ref.Line))
		}

		for _, comment := range WrapString(buff.String(), wrapWidth) {
			if _, err := w.WriteString(fmt.Sprintf("#: %s\n", comment)); err != nil {
				return err
			}
		}
	}

	if len(p.Flags) != 0 {
		if _, err := w.WriteString(fmt.Sprintf("#, %s\n", strings.Join(p.Flags, ", "))); err != nil {
			return err
		}
	}

	return nil
}
