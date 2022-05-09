package spreak

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/text/language"
)

var (
	ErrRequireStringTag = errors.New("spreak: unsupported type, expecting string or language.Tag")

	errMissingLocale = errors.New("spreak: locale missing")
	errSpreak        = errors.New("spreak")
)

// ErrNotFound is the error returned by a loader if no matching context was found.
// If a loader returns any other error, the bundle creation will abort.
type ErrNotFound struct {
	Language   language.Tag
	Type       string
	Identifier string
}

func NewErrNotFound(lang language.Tag, source string, format string, vars ...interface{}) *ErrNotFound {
	return &ErrNotFound{
		Language:   lang,
		Type:       source,
		Identifier: fmt.Sprintf(format, vars...),
	}
}

func (e *ErrNotFound) Error() string { return e.String() }

func (e *ErrNotFound) String() string {
	return fmt.Sprintf("spreak: item of type %q for lang=%q could not be found: %s ", e.Type, e.Language, e.Identifier)
}

// ErrMissingLanguage is the error returned when a Locale should be created and the matching language is not
// loaded or has no Catalog.
type ErrMissingLanguage struct {
	Language language.Tag
}

func newMissingLanguageError(lang language.Tag) *ErrMissingLanguage {
	return &ErrMissingLanguage{Language: lang}
}

func (e *ErrMissingLanguage) Error() string { return e.String() }

func (e *ErrMissingLanguage) String() string {
	return fmt.Sprintf("spreak: language not found: lang=%q ", e.Language)
}

// ErrMissingDomain is the error returned when a domain does not exist for a language.
type ErrMissingDomain struct {
	Language language.Tag
	Domain   string
}

func (e *ErrMissingDomain) Error() string { return e.String() }

func (e *ErrMissingDomain) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("spreak: domain not found: lang=%q ", e.Language))
	if e.Domain != "" {
		b.WriteString(fmt.Sprintf("domain=%q ", e.Domain))
	} else {
		b.WriteString("domain='' (empty string) ")
	}
	return b.String()
}

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
