package po

import (
	"fmt"
	"io"
	"sort"

	"github.com/vorlif/spreak/internal/util"
)

type File struct {
	Header   *Header
	Messages map[string]map[string]*util.Message
}

func (f *File) AddMessage(msg *util.Message) {
	if _, ok := f.Messages[msg.Context]; !ok {
		f.Messages[msg.Context] = make(map[string]*util.Message)
	}

	if _, ok := f.Messages[msg.Context][msg.ID]; ok {
		f.Messages[msg.Context][msg.ID].Merge(msg)
	} else {
		f.Messages[msg.Context][msg.ID] = msg
	}
}

func (f *File) WriteTo(w io.StringWriter, wrapWidth int) error {
	// sort the message as ReferenceFile/ReferenceLine field
	var messages []*util.Message
	for ctx := range f.Messages {
		for _, msg := range f.Messages[ctx] {
			messages = append(messages, msg)
		}
	}
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Comment.Less(messages[j].Comment)
	})

	if f.Header != nil {
		if err := f.Header.WriteTo(w, wrapWidth); err != nil {
			return err
		}

		if _, err := w.WriteString("\n"); err != nil {
			return err
		}
	}
	for _, msg := range messages {
		if err := msg.WriteTo(w, wrapWidth); err != nil {
			return err
		}

		if _, err := w.WriteString("\n"); err != nil {
			return err
		}
	}

	return nil
}

func (f *File) String() string {
	if f.Header != nil {
		return fmt.Sprintf("%s %s", f.Header.ProjectIDVersion, f.Header.Language)
	}
	return "po file"
}
