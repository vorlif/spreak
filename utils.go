package spreak

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"golang.org/x/text/language"
)

var (
	ErrRequireStringTag = errors.New("spreak: unsupported type, expecting string or language.Tag")

	errMissingLocale = errors.New("spreak: locale missing")
)

// ExpandLanguage returns possible filenames for a language tag without extension.
func ExpandLanguage(lang language.Tag) []string {
	expansions := make(map[string]bool, 4)
	expansions[lang.String()] = true

	base, baseConf := lang.Base()
	if baseConf > language.No {
		expansions[base.ISO3()] = true
		expansions[base.String()] = true
	}

	region, regionConf := lang.Region()
	if regionConf > language.No && baseConf > language.No {
		key := fmt.Sprintf("%s_%s", base.String(), region.String())
		expansions[key] = true

		key = fmt.Sprintf("%s-%s", base.String(), region.String())
		expansions[key] = true
	}

	script, scriptConf := lang.Script()
	if scriptConf > language.No && baseConf > language.No {
		key := fmt.Sprintf("%s_%s", base.String(), script.String())
		expansions[key] = true

		key = fmt.Sprintf("%s-%s", base.String(), script.String())
		expansions[key] = true
	}

	return stringMapKeys(expansions)
}

func parseLanguageName(lang string) (language.Tag, error) {
	if idx := strings.Index(lang, ":"); idx != -1 {
		lang = lang[:idx]
	}
	if idx := strings.Index(lang, "@"); idx != -1 {
		lang = lang[:idx]
	}
	if idx := strings.Index(lang, "."); idx != -1 {
		lang = lang[:idx]
	}
	lang = strings.TrimSpace(lang)
	return language.Parse(lang)
}

func languageInterfaceToTag(i interface{}) (language.Tag, error) {
	switch v := i.(type) {
	case string:
		tag, err := parseLanguageName(v)
		if err != nil {
			return language.Und, err
		}
		return tag, nil
	case language.Tag:
		return v, nil
	default:
		return language.Und, ErrRequireStringTag
	}
}

func stringMapKeys(m map[string]bool) []string {
	keys := make([]string, len(m))
	i := 0
	for key := range m {
		keys[i] = key
		i++
	}
	// Longest first
	sort.SliceStable(keys, func(i, j int) bool {
		if x, y := len(keys[i]), len(keys[j]); x != y {
			return x > y
		}

		return keys[i] < keys[j]
	})
	return keys
}

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
