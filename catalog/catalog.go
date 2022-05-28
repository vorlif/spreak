package catalog

import "golang.org/x/text/language"

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

// A Decoder reads and decodes catalogs for a language and a domain from a byte array.
type Decoder interface {
	Decode(lang language.Tag, domain string, data []byte) (Catalog, error)
}
