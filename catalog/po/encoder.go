package po

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/vorlif/spreak/internal/util"
)

// Marshal serializes a Go File as a .po document.
//
// It is a shortcut for Encoder.Encode() with the default options.
func Marshal(f *File) ([]byte, error) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	err := enc.Encode(f)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// An Encoder writes a po file to an output stream.
type Encoder struct {
	w                io.Writer
	wrapWidth        int
	writeHeader      bool
	writeEmptyHeader bool
	writeReferences  bool
	// Function which can be used to sort the messages of a Po file.
	sortFunction func(a *Message, b *Message) int
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w:                w,
		wrapWidth:        -1,
		writeHeader:      true,
		writeReferences:  true,
		writeEmptyHeader: true,
		sortFunction:     DefaultSortFunction,
	}
}

// SetWrapWidth defines at which length the texts should be wrapped.
// To disable wrapping, the value can be set to -1.
// Default is -1.
func (enc *Encoder) SetWrapWidth(wrapWidth int) { enc.wrapWidth = wrapWidth }

// SetWriteHeader sets whether a header should be written or not.
// Default is true.
func (enc *Encoder) SetWriteHeader(write bool) { enc.writeHeader = write }

// SetWriteEmptyHeader sets whether a header without values should also be written or not.
// Default is true.
func (enc *Encoder) SetWriteEmptyHeader(write bool) { enc.writeEmptyHeader = write }

// SetWriteReferences sets whether references to the origin of the text should be stored or not.
// Default is true.
func (enc *Encoder) SetWriteReferences(write bool) { enc.writeReferences = write }

// SetSortFunction can be used to set a function with which the messages can be sorted before writing.
// This sort is not guaranteed to be stable.
// cmp(a, b) should return a negative number when a < b, a positive number when a > b and zero when a == b.
func (enc *Encoder) SetSortFunction(cmp func(a *Message, b *Message) int) { enc.sortFunction = cmp }

func (enc *Encoder) Encode(f *File) error {
	if f == nil {
		return fmt.Errorf("po: cannot encode a nil interface")
	}

	b := enc.encode(f)
	_, err := enc.w.Write(b)
	if err != nil {
		return fmt.Errorf("po: cannot write: %w", err)
	}

	return nil
}

func (enc *Encoder) encode(f *File) []byte {
	// We initialize the buffer with 50Kb, which should be sufficient for medium-sized projects.
	buff := bytes.NewBuffer(make([]byte, 0, 50*1024))

	if f.Header != nil && enc.writeHeader {
		enc.encodeHeader(buff, f.Header)
		buff.WriteString("\n")
	}

	if f.Messages != nil {
		var messages []*Message
		for ctx := range f.Messages {
			for _, msg := range f.Messages[ctx] {
				messages = append(messages, msg)
			}
		}
		if enc.sortFunction != nil {
			slices.SortFunc(messages, enc.sortFunction)
		}
		for i, msg := range messages {
			if i > 0 {
				buff.WriteString("\n")
			}
			enc.encodeMessage(buff, msg)
		}
	}

	return buff.Bytes()
}

type headerEntry struct {
	key   string
	value string
}

func (enc *Encoder) encodeHeader(buff *bytes.Buffer, h *Header) {
	enc.encodeComment(buff, h.Comment)

	headers := []headerEntry{
		{HeaderProjectIDVersion, h.ProjectIDVersion},
		{HeaderReportMsgIDBugsTo, h.ReportMsgidBugsTo},
		{HeaderPOTCreationDate, h.POTCreationDate},
		{HeaderPORevisionDate, h.PORevisionDate},
		{HeaderLastTranslator, h.LastTranslator},
		{HeaderLanguageTeam, h.LanguageTeam},
		{HeaderLanguage, h.Language},
		{HeaderMIMEVersion, h.MimeVersion},
		{HeaderContentType, h.ContentType},
		{HeaderContentTransferEncoding, h.ContentTransferEncoding},
		{HeaderPluralForms, h.PluralForms},
		{HeaderXGenerator, h.XGenerator},
	}

	for k, v := range h.UnknownFields {
		headers = append(headers, headerEntry{k, v})
	}

	if !enc.writeEmptyHeader {
		var hasHeader bool
		for _, header := range headers {
			if header.value != "" {
				hasHeader = true
				break
			}
		}
		if !hasHeader {
			return
		}
	}

	buff.WriteString(`msgid ""
msgstr ""
`)

	for _, header := range headers {
		isOptional := header.key == HeaderMIMEVersion || header.key == HeaderPluralForms || header.key == HeaderXGenerator
		if isOptional && header.value == "" {
			continue
		}
		buff.WriteRune('"')
		buff.WriteString(header.key)
		buff.WriteString(": ")
		buff.WriteString(header.value)
		buff.WriteString("\\n\"\n")
	}
}

func (enc *Encoder) encodeComment(buff *bytes.Buffer, c *Comment) {
	if c == nil {
		return
	}

	if c.Translator != "" {
		for _, comment := range util.WrapString(c.Translator, enc.wrapWidth) {
			buff.WriteString("# ")
			buff.WriteString(comment)
			buff.WriteString("\n")
		}
	}

	if c.Extracted != "" {
		for _, comment := range util.WrapString(c.Extracted, enc.wrapWidth) {
			buff.WriteString("#. ")
			buff.WriteString(comment)
			buff.WriteString("\n")
		}
	}

	if len(c.References) > 0 && enc.writeReferences {
		builder := strings.Builder{}
		for i, ref := range c.References {
			if i > 0 {
				builder.WriteString(" ")
			}

			builder.WriteString(ref.Path)

			if ref.Line > 0 {
				builder.WriteRune(':')
				builder.WriteString(strconv.Itoa(ref.Line))
			}
		}

		for _, comment := range util.WrapString(builder.String(), enc.wrapWidth) {
			buff.WriteString("#: ")
			buff.WriteString(comment)
			buff.WriteString("\n")
		}
	}

	if len(c.Flags) > 0 {
		buff.WriteString("#, ")
		buff.WriteString(strings.Join(c.Flags, ", "))
		buff.WriteString("\n")
	}
}

func (enc *Encoder) encodeMessage(buff *bytes.Buffer, m *Message) {
	enc.encodeComment(buff, m.Comment)

	if m.Context != "" {
		buff.WriteString("msgctxt ")
		buff.WriteString(EncodePoString(m.Context, enc.wrapWidth))
		buff.WriteString("\n")
	}

	buff.WriteString("msgid ")
	buff.WriteString(EncodePoString(m.ID, enc.wrapWidth))
	buff.WriteString("\n")

	hasPlural := m.IDPlural != "" || len(m.Str) > 1
	if hasPlural {
		buff.WriteString("msgid_plural ")
		buff.WriteString(EncodePoString(m.IDPlural, enc.wrapWidth))
		buff.WriteString("\n")
	}

	enc.encodeTranslations(buff, hasPlural, m.Str)
}

func (enc *Encoder) encodeTranslations(buff *bytes.Buffer, plural bool, orig map[int]string) {
	m := make(map[int]string, len(orig))
	for k, v := range orig {
		m[k] = EncodePoString(v, enc.wrapWidth)
	}

	// We need at least one entry
	if len(m) == 0 {
		m[0] = `""`
	}

	if plural {
		if len(m) == 1 {
			// Plural needs at least two entries
			m[1] = `""`
		}

		keys := util.Keys(m)
		sort.Ints(keys)

		for _, k := range keys {
			_, _ = fmt.Fprintf(buff, "msgstr[%d] %s\n", k, m[k])
		}
	} else {
		buff.WriteString("msgstr ")
		buff.WriteString(m[0])
		buff.WriteString("\n")
	}
}

func DefaultSortFunction(a *Message, b *Message) int {
	if a.Comment != nil && b.Comment != nil {
		a.Comment.sort()
		b.Comment.sort()

		for i := 0; i < len(a.Comment.References); i++ {
			if i >= len(b.Comment.References) {
				break
			}
			if c := strings.Compare(a.Comment.References[i].Path, b.Comment.References[i].Path); c != 0 {
				return c
			}
			if c, k := a.Comment.References[i].Line, b.Comment.References[i].Line; c != k {
				if c > k {
					return 1
				} else {
					return -1
				}
			}
			if c, k := a.Comment.References[i].Column, b.Comment.References[i].Column; c != k {
				if c > k {
					return 1
				} else {
					return -1
				}
			}

			if c, k := a.Comment.References[i].Path, b.Comment.References[i].Path; c != "" && k != "" {
				return strings.Compare(c, k)
			}
		}
	}

	if c := strings.Compare(a.Context, b.Context); c != 0 {
		return c
	}

	if c := strings.Compare(a.ID, b.ID); c != 0 {
		return c
	}

	return 0
}
