package catalog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/catalog/cldrplural"
)

// A JSONCatalog represents a collection of translations in JSON format.
type JSONCatalog interface {
	Catalog

	// Domain returns the domain to which this catalog belongs.
	Domain() string
	// Language returns the language to which this catalog belongs.
	Language() language.Tag
	// Messages returns a deep copy of the messages that belong to this catalog.
	Messages() JSONMessages
}

// JSONMessages represents a collection of messages of a JSON catalog.
type JSONMessages map[string]*JSONMessage

// jsonLookupMap is a map for a quick lookup of messages.
// First key is the context and second the MsgID (e.g. lookup["context"]["hello"]).
type jsonLookupMap map[string]map[string]*JSONMessage

type jsonCatalog struct {
	// Language to which this catalog belongs.
	language language.Tag
	// Domain to which this catalog belongs.
	domain string
	// pluralSet is a CLDR rule set for determining the plural category for any number.
	pluralSet *cldrplural.RuleSet
	// Map for a quick lookup of messages.
	// First key is the context and second the msg key (e.g. lookup["context"]["app.name"]).
	lookupMap jsonLookupMap
}

var _ JSONCatalog = (*jsonCatalog)(nil)

// NewJSONCatalog creates a new JSONCatalog for the defined language.
//
// The catalog does not contain any translations, but can be filled with json.Unmarshal from a text file.
func NewJSONCatalog(lang language.Tag, domain string) JSONCatalog {
	pluralSet, _ := cldrplural.ForLanguage(lang)

	cat := &jsonCatalog{
		lookupMap: make(jsonLookupMap),
		domain:    domain,
		language:  lang,
		pluralSet: pluralSet,
	}
	return cat
}

// NewJSONCatalogWithMessages creates a new JSONCatalog for the defined language and messages.
//
// During creation, a deep copy of the messages is created.
// If the messages have an invalid format, an error is returned.
// Plural categories required for the language are added to the messages, and unnecessary ones are removed.
func NewJSONCatalogWithMessages(lang language.Tag, domain string, messages JSONMessages) (JSONCatalog, error) {
	pluralSet, _ := cldrplural.ForLanguage(lang)

	cat := &jsonCatalog{
		lookupMap: make(jsonLookupMap),
		domain:    domain,
		language:  lang,
		pluralSet: pluralSet,
	}

	for key, msg := range messages {
		if err := cat.setMessage(key, msg); err != nil {
			return nil, err
		}
	}

	return cat, nil
}

func (c *jsonCatalog) Lookup(ctx, msgID string) (string, error) {
	tr, err := c.getTranslation(ctx, msgID, cldrplural.Other)
	if err != nil {
		return msgID, err
	}

	return tr, nil
}

func (c *jsonCatalog) LookupPlural(ctx, msgID string, n any) (string, error) {
	cat, errEv := c.pluralSet.Evaluate(n)
	if errEv != nil {
		return msgID, errEv
	}
	tr, err := c.getTranslation(ctx, msgID, cat)
	if err != nil {
		return msgID, err
	}

	return tr, nil
}

// Domain returns the domain to which this catalog belongs.
func (c jsonCatalog) Domain() string { return c.domain }

// Language returns the language to which this catalog belongs.
func (c jsonCatalog) Language() language.Tag { return c.language }

// Messages returns a deep copy of the messages that belong to this catalog.
func (c jsonCatalog) Messages() JSONMessages {
	cpy := make(JSONMessages, len(c.lookupMap))

	for ctx := range c.lookupMap {
		for key := range c.lookupMap[ctx] {
			msg := c.lookupMap[ctx][key]
			cpy[key] = msg.Clone()
		}
	}
	return cpy
}

func (c *jsonCatalog) getTranslation(ctx, key string, cat cldrplural.Category) (string, error) {
	if ctx != "" {
		key += "_" + ctx
	}
	if _, hasCtx := c.lookupMap[ctx]; !hasCtx {
		return "", NewErrMissingContext(c.language, c.domain, ctx)
	}

	if _, hasMsg := c.lookupMap[ctx][key]; !hasMsg {
		return "", NewErrMissingMessageID(c.language, c.domain, ctx, key)
	}

	msg := c.lookupMap[ctx][key]
	tr, ok := msg.Translations[cat]
	if !ok || tr == "" {
		return "", NewErrMissingTranslation(c.language, c.domain, ctx, key, int(cat))
	}

	return tr, nil
}

// Auxiliary method, which adds a message to the lookup map.
// Cannot be made publicly accessible as this would violate the principle of immutability.
func (c *jsonCatalog) setMessage(key string, srcMsg *JSONMessage) error {
	if key == "" {
		return errors.New("spreak: The message key must not be empty")
	}

	if srcMsg == nil {
		return fmt.Errorf("spreak: No message for \"%s\" defined", key)
	}

	if srcMsg.Context != "" && !strings.HasSuffix(key, "_"+srcMsg.Context) {
		return fmt.Errorf("spreak: Key does not match the context - should end with '_%s'", srcMsg.Context)
	}

	msg := srcMsg.Clone()
	if msg.Translations == nil {
		msg.Translations = make(map[cldrplural.Category]string)
	}

	if _, hasOther := msg.Translations[cldrplural.Other]; !hasOther {
		return fmt.Errorf("spreak: \"%s\" does not have an \"other\" value, but is required", key)
	}

	applyCategoriesToJSONMessage(c.pluralSet.Categories, msg)

	ctx := msg.Context
	if _, ok := c.lookupMap[ctx]; !ok {
		c.lookupMap[ctx] = make(map[string]*JSONMessage)
	}

	c.lookupMap[ctx][key] = msg
	return nil
}

// Auxiliary method for tests.
func (c *jsonCatalog) mustSetMessage(msgID string, srcMsg *JSONMessage) {
	if err := c.setMessage(msgID, srcMsg); err != nil {
		panic(err)
	}
}

// UnmarshalJSON can be used to fill the catalog with messages from a JSON catalog file.
//
// The process is not goroutine-safe and should not be performed after the handover to a spreak.Bundle.
func (c *jsonCatalog) UnmarshalJSON(data []byte) error {
	file := make(JSONMessages)
	if err := json.Unmarshal(data, &file); err != nil {
		return err
	}

	if len(file) == 0 {
		return errors.New("spreak: File contains no translations")
	}

	for msgID, msg := range file {
		if err := c.setMessage(msgID, msg); err != nil {
			return err
		}
	}

	return nil
}

func (c jsonCatalog) MarshalJSON() ([]byte, error) {
	messages := make(JSONMessages, len(c.lookupMap))

	for ctx := range c.lookupMap {
		for key := range c.lookupMap[ctx] {
			messages[key] = c.lookupMap[ctx][key]
		}
	}

	keys := slices.Sorted(maps.Keys(messages))

	var buf bytes.Buffer
	buf.WriteRune('{')

	for i, k := range keys {
		if i > 0 {
			buf.WriteRune(',')
		}

		data, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		buf.Write(data)
		buf.WriteRune(':')

		data, err = json.Marshal(messages[k])
		if err != nil {
			return nil, err
		}
		buf.Write(data)
	}
	buf.WriteRune('}')

	return buf.Bytes(), nil
}

// ApplyPluralCategoriesToJSONMessage removes plural categories that do not belong to the language and
// adds those that belong to the language but are still missing.
//
// If only the plural form "Other" is defined, it is assumed that it is a singular entry and
// the plural categories are NOT added.
func ApplyPluralCategoriesToJSONMessage(lang language.Tag, msg *JSONMessage) {
	if msg == nil || msg.Translations == nil {
		return
	}

	pluralSet, _ := cldrplural.ForLanguage(lang)
	applyCategoriesToJSONMessage(pluralSet.Categories, msg)
}

func applyCategoriesToJSONMessage(categories []cldrplural.Category, msg *JSONMessage) {
	// If only "Other" is set, it is a message without a plural and the plural rule does not have to be applied.
	if len(msg.Translations) == 1 {
		if _, hasOther := msg.Translations[cldrplural.Other]; hasOther {
			return
		}
	}

	for cat := range cldrplural.CategoryNames {
		// Remove plural categories that are not supported by this language.
		supported := slices.Contains(categories, cat)
		if !supported {
			delete(msg.Translations, cat)
			continue
		}

		// Add plural categories that are supported by this language and are not yet set.
		if _, hasValue := msg.Translations[cat]; !hasValue {
			msg.Translations[cat] = ""
		}
	}
}

type JSONMessage struct {
	Comment string `json:"comment,omitempty"`
	Context string `json:"context,omitempty"`

	// Translations of this message.
	// For texts without a plural, the map only contains the CLDR category "Other".
	// Otherwise, it contains all plural categories that are required for the language of the catalog.
	Translations map[cldrplural.Category]string
}

// Clone creates a deep copy of the message.
func (m JSONMessage) Clone() *JSONMessage {
	return &JSONMessage{
		Comment:      m.Comment,
		Context:      m.Context,
		Translations: maps.Clone(m.Translations),
	}
}

func (m *JSONMessage) UnmarshalJSON(data []byte) error {
	if m.Translations == nil {
		m.Translations = make(map[cldrplural.Category]string, 1)
	}

	// Check whether only the CLDR "Other" category is set.
	var other string
	if err := json.Unmarshal(data, &other); err == nil {
		m.Translations[cldrplural.Other] = other
		return nil
	}

	// Check whether several values are set.
	var mm map[string]string
	if err := json.Unmarshal(data, &mm); err != nil {
		return err
	}

	for key, value := range mm {
		switch strings.ToLower(key) {
		case "comment":
			m.Comment = value
		case "context":
			m.Context = value
		case "zero":
			m.Translations[cldrplural.Zero] = value
		case "one":
			m.Translations[cldrplural.One] = value
		case "two":
			m.Translations[cldrplural.Two] = value
		case "few":
			m.Translations[cldrplural.Few] = value
		case "many":
			m.Translations[cldrplural.Many] = value
		case "other":
			m.Translations[cldrplural.Other] = value
		}
	}

	return nil
}

// Defines the order in which the data is to be written.
var categories = []cldrplural.Category{
	cldrplural.Zero,
	cldrplural.One,
	cldrplural.Two,
	cldrplural.Few,
	cldrplural.Many,
	cldrplural.Other,
}

func (m JSONMessage) MarshalJSON() ([]byte, error) {
	if m.Translations == nil {
		m.Translations = make(map[cldrplural.Category]string)
	}

	if len(m.Translations) == 0 {
		m.Translations[cldrplural.Other] = ""
	}

	// If only Other is set, only Other is returned.
	if m.Comment == "" && m.Context == "" && len(m.Translations) == 1 {
		other := m.Translations[cldrplural.Other]
		return json.Marshal(other)
	}

	var buf bytes.Buffer

	writeKeyValue := func(key string, value any) error {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		buf.WriteString(fmt.Sprintf(`"%s": `, key))
		buf.Write(data)
		return nil
	}

	buf.WriteRune('{')
	if m.Comment != "" {
		if err := writeKeyValue("comment", m.Comment); err != nil {
			return nil, err
		}
		buf.WriteRune(',')
	}

	if m.Context != "" {
		if err := writeKeyValue("context", m.Context); err != nil {
			return nil, err
		}
		buf.WriteRune(',')
	}

	i := 0
	for _, cat := range categories {
		value, ok := m.Translations[cat]
		if !ok {
			continue
		}

		if i > 0 {
			buf.WriteRune(',')
		}
		i++

		key := strings.ToLower(cat.String())
		if err := writeKeyValue(key, value); err != nil {
			return nil, err
		}
	}

	buf.WriteRune('}')
	return buf.Bytes(), nil
}
