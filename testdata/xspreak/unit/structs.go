package unit

import "github.com/vorlif/spreak/localize"

type TranslationStruct struct {
	Singular localize.Singular
	plural   localize.Plural
	Domain   localize.Domain
	context  localize.Context
}
