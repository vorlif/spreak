package catalog

import (
	"fmt"
	"strings"

	"golang.org/x/text/language"
)

// ErrMissingContext is the error returned when a matching context was not found for a language and domain.
type ErrMissingContext struct {
	Language language.Tag
	Domain   string
	Context  string
}

func (e *ErrMissingContext) Error() string { return e.String() }

func (e *ErrMissingContext) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("spreak: context not found: lang=%q ", e.Language))
	if e.Domain != "" {
		b.WriteString(fmt.Sprintf("domain=%q ", e.Domain))
	}
	if e.Context != "" {
		b.WriteString(fmt.Sprintf("ctx=%q ", e.Context))
	} else {
		b.WriteString("ctx='' (empty string)")
	}
	return b.String()
}

// ErrMissingMessageID is the error returned when a matching message was not found for a language and domain.
type ErrMissingMessageID struct {
	Language language.Tag
	Domain   string
	Context  string
	MsgID    string
}

func (e *ErrMissingMessageID) Error() string { return e.String() }

func (e *ErrMissingMessageID) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("spreak: msgID not found: lang=%q ", e.Language))
	if e.Domain != "" {
		b.WriteString(fmt.Sprintf("domain=%q ", e.Domain))
	}
	if e.Context != "" {
		b.WriteString(fmt.Sprintf("ctx=%q ", e.Context))
	}
	b.WriteString(fmt.Sprintf("msgID=%q", e.MsgID))
	return b.String()
}

// ErrMissingTranslation is the error returned when there is no translation for a domain of a language for a message.
type ErrMissingTranslation struct {
	Language language.Tag
	Domain   string
	Context  string
	MsgID    string
	Idx      int
}

func (e *ErrMissingTranslation) Error() string { return e.String() }

func (e *ErrMissingTranslation) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("spreak: translation not found: lang=%q ", e.Language))
	if e.Domain != "" {
		b.WriteString(fmt.Sprintf("domain=%q ", e.Domain))
	}
	if e.Context != "" {
		b.WriteString(fmt.Sprintf("ctx=%q ", e.Context))
	}
	b.WriteString(fmt.Sprintf("msgID=%q idx=%d", e.MsgID, e.Idx))
	return b.String()
}
