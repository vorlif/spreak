package spreak

import (
	"golang.org/x/text/language"
)

// Catalog represents a collection of messages (translations) for a language and a domain.
// Normally it is a PO or MO file.
type Catalog interface {
	// GetTranslation Returns a translation for an ID within a given context.
	GetTranslation(ctx, msgID string) (string, error)
	// GetPluralTranslation Returns a translation within a given context.
	// Here n is a number that should be used to determine the plural form.
	GetPluralTranslation(ctx, msgID string, n interface{}) (string, error)

	Language() language.Tag
}

type gettextCatalog struct {
	language language.Tag

	translations messageLookupMap
	pluralFunc   pluralFunction
	domain       string
}

type gettextMessage struct {
	Context      string
	ID           string
	IDPlural     string
	Translations map[int]string
}

// Map for a quick lookup of messages.
// First key is the context and second the MsgID (e.g. lookup["context"]["hello"]).
type messageLookupMap map[string]map[string]*gettextMessage

var _ Catalog = (*gettextCatalog)(nil)

func (c *gettextCatalog) GetTranslation(ctx, msgID string) (string, error) {
	msg, err := c.getMessage(ctx, msgID, 0)
	if err != nil {
		return msgID, err
	}

	return msg.Translations[0], nil
}

func (c *gettextCatalog) GetPluralTranslation(ctx, msgID string, n interface{}) (string, error) {
	idx := c.pluralFunc(n)
	msg, err := c.getMessage(ctx, msgID, idx)
	if err != nil {
		if idx == 0 {
			return msgID, err
		}

		return msgID, err
	}

	return msg.Translations[idx], nil
}

func (c *gettextCatalog) Language() language.Tag { return c.language }

func (c *gettextCatalog) getMessage(ctx, msgID string, idx int) (*gettextMessage, error) {
	if _, hasCtx := c.translations[ctx]; !hasCtx {
		err := &ErrMissingContext{
			Language: c.language,
			Domain:   c.domain,
			Context:  ctx,
		}
		return nil, err
	}

	if _, hasMsg := c.translations[ctx][msgID]; !hasMsg {
		err := &ErrMissingMessageID{
			Language: c.language,
			Domain:   c.domain,
			Context:  ctx,
			MsgID:    msgID,
		}
		return nil, err
	}

	msg := c.translations[ctx][msgID]
	if tr, hasTranslation := msg.Translations[idx]; !hasTranslation || tr == "" {
		err := &ErrMissingTranslation{
			Language: c.language,
			Domain:   c.domain,
			Context:  ctx,
			MsgID:    msgID,
			Idx:      idx,
		}
		return nil, err
	}

	return msg, nil
}
