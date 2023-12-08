package po

import (
	"fmt"
	"strings"
	"time"
)

const (
	HeaderProjectIDVersion        = "Project-Id-Version"
	HeaderReportMsgIDBugsTo       = "Report-Msgid-Bugs-To"
	HeaderPOTCreationDate         = "POT-Creation-Date"
	HeaderPORevisionDate          = "PO-Revision-Date"
	HeaderLastTranslator          = "Last-Translator"
	HeaderLanguageTeam            = "Language-Team"
	HeaderLanguage                = "Language"
	HeaderMIMEVersion             = "MIME-Version"
	HeaderContentType             = "Content-Type"
	HeaderContentTransferEncoding = "Content-Transfer-Encoding"
	HeaderPluralForms             = "Plural-Forms"
	HeaderXGenerator              = "X-Generator"
)

// File is Go data structure for a PO file.
type File struct {
	Header   *Header
	Messages Messages
}

// Messages are the messages of a Po file.
// The first map references the context, the second map the ID of the message.
type Messages map[string]map[string]*Message

// Add a message to the File.
// If there is already a message with the same ID, the comments are merged.
func (m Messages) Add(msg *Message) {
	if _, ok := m[msg.Context]; !ok {
		m[msg.Context] = make(map[string]*Message)
	}

	if _, ok := m[msg.Context][msg.ID]; ok {
		m[msg.Context][msg.ID].Merge(msg)
	} else {
		m[msg.Context][msg.ID] = msg
	}
}

func NewFile() *File {
	return &File{
		Header:   &Header{},
		Messages: make(Messages),
	}
}

// AddMessage adds a message to the File.
// If there is already a message with the same ID, the comments are merged.
func (f *File) AddMessage(msg *Message) {
	if f.Messages == nil {
		f.Messages = make(Messages)
	}

	if msg.ID == "" {
		return
	}

	f.Messages.Add(msg)
}

// GetMessage returns the message for a context and an ID.
// If no ID exists, nil is returned.
func (f *File) GetMessage(ctx string, id string) *Message {
	if _, hasCtx := f.Messages[ctx]; !hasCtx {
		return nil
	}

	msg, ok := f.Messages[ctx][id]
	if !ok {
		return nil
	}

	return msg
}

func (f *File) String() string {
	if f.Header != nil {
		return fmt.Sprintf("po file %s %s", f.Header.ProjectIDVersion, f.Header.Language)
	}
	return "po file"
}

type Header struct {
	Comment                 *Comment // Header Comments
	ProjectIDVersion        string   // Project-Id-Version: PACKAGE VERSION
	ReportMsgidBugsTo       string   // Report-Msgid-Bugs-To: FIRST AUTHOR <EMAIL@ADDRESS>
	POTCreationDate         string   // POT-Creation-Date: YEAR-MO-DA HO:MI+ZONE
	PORevisionDate          string   // PO-Revision-Date: YEAR-MO-DA HO:MI+ZONE
	LastTranslator          string   // Last-Translator: FIRST AUTHOR <EMAIL@ADDRESS>
	LanguageTeam            string   // Language-Team:
	Language                string   // Language: de
	MimeVersion             string   // MIME-Version: 1.0
	ContentType             string   // Content-Type: text/plain; charset=UTF-8
	ContentTransferEncoding string   // Content-Transfer-Encoding: 8bit
	PluralForms             string   // Plural-Forms: nplurals=2; plural=(n != 1);
	XGenerator              string   // X-Generator: Poedit 3.0.1
	UnknownFields           map[string]string
}

// Deprecated: Use Set.
func (h *Header) SetField(key, val string) { h.Set(key, val) }

// Set tores the key and value in the header.
// If a key is already stored, the value is overwritten.
func (h *Header) Set(key, val string) {
	switch strings.ToUpper(key) {
	case strings.ToUpper(HeaderProjectIDVersion):
		h.ProjectIDVersion = val
	case strings.ToUpper(HeaderReportMsgIDBugsTo):
		h.ReportMsgidBugsTo = val
	case strings.ToUpper(HeaderPOTCreationDate):
		h.POTCreationDate = val
	case strings.ToUpper(HeaderPORevisionDate):
		h.PORevisionDate = val
	case strings.ToUpper(HeaderLastTranslator):
		h.LastTranslator = val
	case strings.ToUpper(HeaderLanguageTeam):
		h.LanguageTeam = val
	case strings.ToUpper(HeaderLanguage):
		h.Language = val
	case strings.ToUpper(HeaderMIMEVersion):
		h.MimeVersion = val
	case strings.ToUpper(HeaderContentType):
		h.ContentType = val
	case strings.ToUpper(HeaderContentTransferEncoding):
		h.ContentTransferEncoding = val
	case strings.ToUpper(HeaderPluralForms):
		h.PluralForms = val
	case strings.ToUpper(HeaderXGenerator):
		h.XGenerator = val
	default:
		if h.UnknownFields == nil {
			h.UnknownFields = make(map[string]string)
		}
		h.UnknownFields[key] = val
	}
}

// Get returns the respective value for a key.
// If the key is not contained in the header, an empty string is returned.
func (h *Header) Get(key string) string {
	switch strings.ToUpper(key) {
	case strings.ToUpper(HeaderProjectIDVersion):
		return h.ProjectIDVersion
	case strings.ToUpper(HeaderReportMsgIDBugsTo):
		return h.ReportMsgidBugsTo
	case strings.ToUpper(HeaderPOTCreationDate):
		return h.POTCreationDate
	case strings.ToUpper(HeaderPORevisionDate):
		return h.PORevisionDate
	case strings.ToUpper(HeaderLastTranslator):
		return h.LastTranslator
	case strings.ToUpper(HeaderLanguageTeam):
		return h.LanguageTeam
	case strings.ToUpper(HeaderLanguage):
		return h.Language
	case strings.ToUpper(HeaderMIMEVersion):
		return h.MimeVersion
	case strings.ToUpper(HeaderContentType):
		return h.ContentType
	case strings.ToUpper(HeaderContentTransferEncoding):
		return h.ContentTransferEncoding
	case strings.ToUpper(HeaderPluralForms):
		return h.PluralForms
	case strings.ToUpper(HeaderXGenerator):
		return h.XGenerator
	}

	if h.UnknownFields == nil {
		return ""
	}

	key = strings.ToUpper(key)
	for unknownHeader, val := range h.UnknownFields {
		if strings.ToUpper(unknownHeader) == key {
			return val
		}
	}

	return ""
}

// PlaceholderHeader creates a placeholder header for an empty PO file.
func PlaceholderHeader(packageName, copyrightHolder, bugsAddress string) *Header {
	headerComment := fmt.Sprintf(`SOME DESCRIPTIVE TITLE.
Copyright (C) YEAR %s
This file is distributed under the same license as the %s package.
FIRST AUTHOR <EMAIL@ADDRESS>, YEAR.
`, copyrightHolder, packageName)
	return &Header{
		Comment: &Comment{
			Translator:     headerComment,
			Extracted:      "",
			References:     nil,
			Flags:          []string{"fuzzy"},
			PrevMsgContext: "",
			PrevMsgID:      "",
		},
		ProjectIDVersion:        packageName,
		ReportMsgidBugsTo:       bugsAddress,
		POTCreationDate:         time.Now().Format("2006-01-02 15:04-0700"),
		PORevisionDate:          "YEAR-MO-DA HO:MI+ZONE",
		LastTranslator:          "FULL NAME <EMAIL@ADDRESS>",
		LanguageTeam:            "LANGUAGE <LL@li.org>",
		Language:                "",
		MimeVersion:             "1.0",
		ContentType:             "text/plain; charset=UTF-8",
		ContentTransferEncoding: "8bit",
		PluralForms:             "", // alternative  "nplurals=INTEGER; plural=EXPRESSION;"
	}
}
