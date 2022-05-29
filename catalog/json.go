package catalog

import (
	"encoding/json"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/internal/cldrplural"
)

type jsonDecoder struct{}

func NewJSONDecoder() Decoder { return jsonDecoder{} }

func (jsonDecoder) Decode(lang language.Tag, domain string, data []byte) (Catalog, error) {
	var messages JSONFile
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, err
	}

	catl := &JSONCatalog{
		lookupMap: make(map[string]map[string]*JSONMessage),
		domain:    domain,
		language:  lang,
	}

	catl.pluralSet, _ = cldrplural.ForLanguage(lang)

	for key, msg := range messages {
		if key == "" || msg == nil || msg.Other == "" {
			continue
		}

		if _, ok := catl.lookupMap[msg.Context]; !ok {
			catl.lookupMap[msg.Context] = make(map[string]*JSONMessage)
		}

		catl.lookupMap[msg.Context][key] = msg
	}

	return catl, nil
}

type JSONCatalog struct {
	// Map for a quick lookup of messages.
	// First key is the context and second the msg key (e.g. lookup["context"]["app.name"]).
	lookupMap map[string]map[string]*JSONMessage
	domain    string
	language  language.Tag
	pluralSet *cldrplural.RuleSet
}

func (m *JSONCatalog) GetTranslation(ctx, msgID string) (string, error) {
	tr, err := m.getTranslation(ctx, msgID, cldrplural.Other)
	if err != nil {
		return msgID, err
	}

	return tr, nil
}

func (m *JSONCatalog) GetPluralTranslation(ctx, msgID string, n interface{}) (string, error) {
	cat := m.pluralSet.Evaluate(n)
	tr, err := m.getTranslation(ctx, msgID, cat)
	if err != nil {
		return msgID, err
	}

	return tr, nil
}

func (m JSONCatalog) Language() language.Tag { return m.language }

func (m *JSONCatalog) getTranslation(ctx, key string, cat cldrplural.Category) (string, error) {
	if ctx != "" {
		key += "_" + ctx
	}
	if _, hasCtx := m.lookupMap[ctx]; !hasCtx {
		err := &ErrMissingContext{
			Language: m.language,
			Domain:   m.domain,
			Context:  ctx,
		}
		return "", err
	}

	if _, hasMsg := m.lookupMap[ctx][key]; !hasMsg {
		err := &ErrMissingMessageID{
			Language: m.language,
			Domain:   m.domain,
			Context:  ctx,
			MsgID:    key,
		}
		return "", err
	}

	msg := m.lookupMap[ctx][key]
	tr := msg.getTranslation(cat)
	if tr == "" {
		err := &ErrMissingTranslation{
			Language: m.language,
			Domain:   m.domain,
			Context:  ctx,
			MsgID:    key,
			Idx:      int(cat),
		}
		return "", err
	}

	return tr, nil
}

type JSONFile map[string]*JSONMessage

type jsonMessageAlias JSONMessage

type JSONMessage struct {
	Comment string `json:"comment,omitempty"`
	Context string `json:"context,omitempty"`

	Zero  string `json:"zero,omitempty"`
	One   string `json:"one,omitempty"`
	Two   string `json:"two,omitempty"`
	Few   string `json:"few,omitempty"`
	Many  string `json:"many,omitempty"`
	Other string `json:"other"`
}

func (m *JSONMessage) getTranslation(cat cldrplural.Category) string {
	switch cat {
	case cldrplural.Zero:
		return m.Zero
	case cldrplural.One:
		return m.One
	case cldrplural.Two:
		return m.Two
	case cldrplural.Few:
		return m.Few
	case cldrplural.Many:
		return m.Many
	default:
		return m.Other
	}
}

func (m *JSONMessage) hasFilledFields() bool {
	if m.Comment != "" {
		return true
	}
	if m.Context != "" {
		return true
	}
	if m.Zero != "" {
		return true
	}
	if m.One != "" {
		return true
	}
	if m.Two != "" {
		return true
	}
	if m.Few != "" {
		return true
	}
	if m.Many != "" {
		return true
	}
	return false
}

func (m *JSONMessage) MarshalJSON() ([]byte, error) {
	if m.hasFilledFields() {
		return json.Marshal(*m)
	}

	return json.Marshal(m.Other)
}

func (m *JSONMessage) UnmarshalJSON(data []byte) error {
	var other string
	if err := json.Unmarshal(data, &other); err == nil {
		m.Other = other
		return nil
	}

	aux := &struct{ *jsonMessageAlias }{jsonMessageAlias: (*jsonMessageAlias)(m)}
	return json.Unmarshal(data, aux)
}
