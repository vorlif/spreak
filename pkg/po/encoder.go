package po

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/vorlif/spreak/internal/util"
)

// An Encoder writes a po file to an output stream.
type Encoder struct {
	w               *bufio.Writer
	wrapWidth       int
	writeHeader     bool
	writeReferences bool
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w:               bufio.NewWriter(w),
		wrapWidth:       -1,
		writeHeader:     true,
		writeReferences: true,
	}
}

func (e *Encoder) SetWrapWidth(wrapWidth int) {
	e.wrapWidth = wrapWidth
}

func (e *Encoder) SetWriteHeader(write bool) {
	e.writeHeader = write
}

func (e *Encoder) SetWriteReferences(write bool) {
	e.writeReferences = write
}

func (e *Encoder) Encode(f *File) error {
	if f.Header != nil && e.writeHeader {
		if err := e.encodeHeader(f.Header); err != nil {
			return err
		}

		if _, err := e.w.WriteString("\n"); err != nil {
			return err
		}
	}

	if f.Messages != nil {
		var messages []*Message
		for ctx := range f.Messages {
			for _, msg := range f.Messages[ctx] {
				messages = append(messages, msg)
			}
		}
		sort.Slice(messages, func(i, j int) bool {
			return messages[i].Comment.Less(messages[j].Comment)
		})
		for _, msg := range messages {
			if err := e.encodeMessage(msg); err != nil {
				return err
			}

			if _, err := e.w.WriteString("\n"); err != nil {
				return err
			}
		}
	}

	return e.w.Flush()
}

func (e *Encoder) encodeHeader(h *Header) error {
	if err := e.encodeComment(h.Comment); err != nil {
		return err
	}

	lines := []string{
		`msgid ""` + "\n",
		`msgstr ""` + "\n",
		fmt.Sprintf(`"%s: %s\n"`+"\n", HeaderProjectIDVersion, h.ProjectIDVersion),
		fmt.Sprintf(`"%s: %s\n"`+"\n", HeaderReportMsgIDBugsTo, h.ReportMsgidBugsTo),
		fmt.Sprintf(`"%s: %s\n"`+"\n", HeaderPOTCreationDate, h.POTCreationDate),
		fmt.Sprintf(`"%s: %s\n"`+"\n", HeaderPORevisionDate, h.PORevisionDate),
		fmt.Sprintf(`"%s: %s\n"`+"\n", HeaderLastTranslator, h.LastTranslator),
		fmt.Sprintf(`"%s: %s\n"`+"\n", HeaderLanguageTeam, h.LanguageTeam),
		fmt.Sprintf(`"%s: %s\n"`+"\n", HeaderLanguage, h.Language),
	}

	if h.MimeVersion != "" {
		lines = append(lines, fmt.Sprintf(`"%s: %s\n"`+"\n", HeaderMIMEVersion, h.MimeVersion))
	}

	lines = append(lines,
		fmt.Sprintf(`"%s: %s\n"`+"\n", HeaderContentType, h.ContentType),
		fmt.Sprintf(`"%s: %s\n"`+"\n", HeaderContentTransferEncoding, h.ContentTransferEncoding),
	)

	if h.PluralForms != "" {
		lines = append(lines, fmt.Sprintf(`"%s: %s\n"`+"\n", HeaderPluralForms, h.PluralForms))
	}

	if h.XGenerator != "" {
		lines = append(lines, fmt.Sprintf(`"%s: %s\n"`+"\n", HeaderXGenerator, h.XGenerator))
	}
	for k, v := range h.UnknownFields {
		lines = append(lines, fmt.Sprintf(`"%s: %s\n"`+"\n", k, v))
	}

	for _, line := range lines {
		if _, err := e.w.WriteString(line); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeComment(c *Comment) error {
	if c == nil {
		return nil
	}

	if c.Translator != "" {
		for _, comment := range util.WrapString(c.Translator, e.wrapWidth) {
			if _, err := e.w.WriteString(fmt.Sprintf("# %s\n", comment)); err != nil {
				return err
			}
		}
	}

	if c.Extracted != "" {
		for _, comment := range util.WrapString(c.Extracted, e.wrapWidth) {
			if _, err := e.w.WriteString(fmt.Sprintf("#. %s\n", comment)); err != nil {
				return err
			}
		}
	}

	if len(c.References) > 0 && e.writeReferences {
		var buff bytes.Buffer
		for _, ref := range c.References {
			buff.WriteString(fmt.Sprintf("%s:%d ", ref.Path, ref.Line))
		}

		for _, comment := range util.WrapString(buff.String(), e.wrapWidth) {
			if _, err := e.w.WriteString(fmt.Sprintf("#: %s\n", comment)); err != nil {
				return err
			}
		}
	}

	if len(c.Flags) > 0 {
		if _, err := e.w.WriteString(fmt.Sprintf("#, %s\n", strings.Join(c.Flags, ", "))); err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) encodeMessage(m *Message) error {
	if err := e.encodeComment(m.Comment); err != nil {
		return err
	}

	if m.Context != "" {
		ctx := fmt.Sprintf("msgctxt %s\n", EncodePoString(m.Context, e.wrapWidth))
		if _, err := e.w.WriteString(ctx); err != nil {
			return err
		}
	}

	msgID := fmt.Sprintf("msgid %s\n", EncodePoString(m.ID, e.wrapWidth))
	if _, err := e.w.WriteString(msgID); err != nil {
		return err
	}

	hasPlural := m.IDPlural != "" || len(m.Str) > 1
	if hasPlural {
		pluralID := fmt.Sprintf("msgid_plural %s\n", EncodePoString(m.IDPlural, e.wrapWidth))
		if _, err := e.w.WriteString(pluralID); err != nil {
			return err
		}
	}

	if err := e.encodeTranslations(hasPlural, m.Str); err != nil {
		return err
	}

	return nil
}

func (e *Encoder) encodeTranslations(plural bool, orig map[int]string) error {
	m := make(map[int]string, len(orig))
	for k, v := range orig {
		m[k] = v
	}

	// We need at least one entry
	if len(m) == 0 {
		m[0] = ""
	}

	if plural {
		if len(m) == 1 {
			// Plural needs at least two entries
			m[1] = ""
		}

		keys := make([]int, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		for _, k := range keys {
			if _, err := e.w.WriteString(fmt.Sprintf("msgstr[%d] \"%s\"\n", k, m[k])); err != nil {
				return err
			}
		}
	} else {
		if _, err := e.w.WriteString(fmt.Sprintf("msgstr \"%s\"\n", m[0])); err != nil {
			return err
		}
	}

	return nil
}
