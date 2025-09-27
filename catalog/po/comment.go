package po

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

type Reference struct {
	Path   string
	Line   int
	Column int
}

func (r Reference) String() string {
	if r.Line > 0 {
		return fmt.Sprintf("%s:%d:%d", r.Path, r.Line, r.Column)
	}

	return r.Path
}

func (r Reference) Equal(o *Reference) bool {
	return r.Path == o.Path && r.Line == o.Line && r.Column == o.Column
}

type Comment struct {
	Translator     string       // #  translator-comments
	Extracted      string       // #. extracted-comments
	References     []*Reference // #: src/file.go:210
	Flags          []string     // #, fuzzy,go-format,range:0..10
	PrevMsgContext string       // #| msgctxt previous-context
	PrevMsgID      string       // #| msgid previous-untranslated-string
}

func NewComment() *Comment {
	return &Comment{
		References: []*Reference{},
		Flags:      []string{},
	}
}

func (c *Comment) AddReference(ref *Reference) {
	if len(c.References) == 0 {
		c.References = []*Reference{ref}
		return
	}

	c.References = append(c.References, ref)
	c.sortReferences()
}

func (c *Comment) HasFlag(flag string) bool {
	return slices.Contains(c.Flags, flag)
}

func (c *Comment) AddFlag(flag string) {
	if c.HasFlag(flag) {
		return
	}
	c.Flags = append(c.Flags, flag)
	c.sortFlags()
}

func (c *Comment) Merge(other *Comment) {
	if other == nil {
		return
	}

	if other.Translator != "" {
		left := strings.Split(c.Translator, "\n")
		right := strings.Split(other.Translator, "\n")
		res := mergeStringArrays(left, right)
		c.Translator = strings.TrimSpace(strings.Join(res, "\n"))
	}

	if other.Extracted != "" {
		left := strings.Split(c.Extracted, "\n")
		right := strings.Split(other.Extracted, "\n")
		res := mergeStringArrays(left, right)
		c.Extracted = strings.TrimSpace(strings.Join(res, "\n"))
	}

	c.Flags = mergeStringArrays(c.Flags, other.Flags)
	c.mergeReferences(other)

	c.sort()
}

func (c *Comment) mergeReferences(other *Comment) {
	for _, otherRef := range other.References {
		if slices.ContainsFunc(c.References, func(ref *Reference) bool { return ref.Equal(otherRef) }) {
			continue
		}

		c.References = append(c.References, otherRef)
	}
}

func (c *Comment) sort() {
	c.sortFlags()
	c.sortReferences()
}

func (c *Comment) sortFlags() {
	slices.Sort(c.Flags)
}

func (c *Comment) sortReferences() {
	sort.Slice(c.References, func(i, j int) bool {
		if c.References[i].Path != c.References[j].Path {
			return c.References[i].Path < c.References[j].Path
		}

		if c.References[i].Line != c.References[j].Line {
			return c.References[i].Line < c.References[j].Line
		}

		return c.References[i].Column < c.References[j].Column
	})
}

func mergeStringArrays(left, right []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(left)+len(right))

	for _, item := range left {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	for _, item := range right {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}
