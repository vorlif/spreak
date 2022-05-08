package util

import (
	"fmt"
	"io"
)

type Message struct {
	Comment  *Comment
	Context  string
	ID       string
	IDPlural string
	Str      map[int]string
}

func NewMessage() *Message {
	return &Message{
		Comment: &Comment{
			References: make([]*Reference, 0),
			Flags:      make([]string, 0),
		},
		Context:  "",
		ID:       "",
		IDPlural: "",
		Str:      nil,
	}
}

func (p *Message) AddReference(ref *Reference) {
	if p.Comment == nil {
		p.Comment = &Comment{}
	}

	p.Comment.AddReference(ref)
}

func (p *Message) WriteTo(w io.StringWriter, wrapWidth int) error {

	if err := p.Comment.WriteTo(w, wrapWidth); err != nil {
		return err
	}

	if p.Context != "" {
		ctx := fmt.Sprintf("msgctxt %s\n", EncodePoString(p.Context, wrapWidth))
		if _, err := w.WriteString(ctx); err != nil {
			return err
		}
	}

	msgID := fmt.Sprintf("msgid %s\n", EncodePoString(p.ID, wrapWidth))
	if _, err := w.WriteString(msgID); err != nil {
		return err
	}

	if p.IDPlural != "" {
		pluralID := fmt.Sprintf("msgid_plural %s\n", EncodePoString(p.IDPlural, wrapWidth))
		if _, err := w.WriteString(pluralID); err != nil {
			return err
		}

		if _, err := w.WriteString(`msgstr[0] ""` + "\n"); err != nil {
			return err
		}
		if _, err := w.WriteString(`msgstr[1] ""` + "\n"); err != nil {
			return err
		}
	} else {
		if _, err := w.WriteString(`msgstr ""` + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func (p *Message) Less(q *Message) bool {
	if p.Comment.Less(q.Comment) {
		return true
	}

	if a, b := p.Context, q.Context; a != b {
		return a < b
	}
	if a, b := p.ID, q.ID; a != b {
		return a < b
	}
	if a, b := p.IDPlural, q.IDPlural; a != b {
		return a < b
	}
	return false
}

func (p *Message) Merge(other *Message) {
	newReferences := make([]*Reference, 0)
	for _, ref := range p.Comment.References {
		for _, otherRef := range other.Comment.References {
			if ref.Line == otherRef.Line && ref.Path == otherRef.Path {
				continue
			}

			newReferences = append(newReferences, otherRef)
		}
	}
	p.Comment.References = append(p.Comment.References, newReferences...)

	newFlags := make([]string, 0)
	for _, flag := range p.Comment.Flags {
		for _, otherFlag := range other.Comment.Flags {
			if flag == otherFlag {
				continue
			}

			newFlags = append(newFlags, otherFlag)
		}
	}
	p.Comment.Flags = append(p.Comment.Flags, newFlags...)

	if p.IDPlural == "" && other.IDPlural != "" {
		p.IDPlural = other.IDPlural
	}
}
