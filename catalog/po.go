package catalog

import (
	"maps"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/catalog/poplural"
)

// PoLookupMap is a map for a quick lookup of messages.
// First key is the context and second the MsgID (e.g. lookup["context"]["hello"]).
type PoLookupMap map[string]map[string]*GettextMessage

type GettextCatalog struct {
	// Language to which this catalog belongs.
	language language.Tag
	// Domain to which this catalog belongs.
	domain string
	// pluralFunc is function for determining the plural form for any number.
	pluralFunc poplural.PluralFunc
	// Map for quick access to the translations.
	lookupMap PoLookupMap
}

var _ Catalog = (*GettextCatalog)(nil)

func (c *GettextCatalog) Lookup(ctx, msgID string) (string, error) {
	msg, err := c.getMessage(ctx, msgID, 0)
	if err != nil {
		return msgID, err
	}

	return msg.translations[0], nil
}

func (c *GettextCatalog) LookupPlural(ctx, msgID string, n any) (string, error) {
	idx, errPlural := c.pluralFunc(n)
	if errPlural != nil {
		return msgID, errPlural
	}
	msg, err := c.getMessage(ctx, msgID, idx)
	if err != nil {
		return msgID, err
	}

	return msg.translations[idx], nil
}

// Domain returns the domain to which this catalog belongs.
func (c GettextCatalog) Domain() string { return c.domain }

// Language returns the language to which this catalog belongs.
func (c GettextCatalog) Language() language.Tag { return c.language }

// Messages returns a deep copy of the messages that belong to this catalog.
// The messages are returned as a nested map, with the first map containing the context
// and the second map containing the message ID.
// Access is therefore via message["context"]["msgId"].
func (c *GettextCatalog) Messages() PoLookupMap {
	cpy := make(PoLookupMap, len(c.lookupMap))
	for ctx := range c.lookupMap {
		cpy[ctx] = make(map[string]*GettextMessage, len(c.lookupMap[ctx]))

		for msgId := range c.lookupMap[ctx] {
			msg := cpy[ctx][msgId]
			cpy[ctx][msgId] = msg.Clone()
		}
	}
	return cpy
}

func (c *GettextCatalog) getMessage(ctx, msgID string, idx int) (*GettextMessage, error) {
	if _, hasCtx := c.lookupMap[ctx]; !hasCtx {
		return nil, NewErrMissingContext(c.language, c.domain, ctx)
	}

	if _, hasMsg := c.lookupMap[ctx][msgID]; !hasMsg {
		return nil, NewErrMissingMessageID(c.language, c.domain, ctx, msgID)
	}

	msg := c.lookupMap[ctx][msgID]
	if tr, hasTranslation := msg.translations[idx]; !hasTranslation || tr == "" {
		return nil, NewErrMissingTranslation(c.language, c.domain, ctx, msgID, idx)
	}

	return msg, nil
}

type GettextMessage struct {
	translations map[int]string
}

// Translations returns a depth copy of the translations.
// The structure is returned as a map, with the key addressing the plural form.
// If there is no plural form, the map contains only one item for 0.
func (m GettextMessage) Translations() map[int]string { return maps.Clone(m.translations) }

// Clone creates a deep copy of the message.
func (m GettextMessage) Clone() *GettextMessage {
	return &GettextMessage{translations: maps.Clone(m.translations)}
}
