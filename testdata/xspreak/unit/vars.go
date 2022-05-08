package unit

import "github.com/vorlif/spreak/localize"

var (
	one         localize.Singular = "one-singular"
	two         localize.Singular = "two-singular"
	nonLocalize                   = "non-localize"
)

var three localize.Singular = "three-singular"

var ignored = "ignored"

var pluralIgnored localize.Plural = "plural without singular is ignored"
