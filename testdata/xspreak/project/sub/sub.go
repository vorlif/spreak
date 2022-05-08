package sub

import (
	"github.com/vorlif/testdata/foo"

	"github.com/vorlif/spreak/localize"
)

type Sub struct {
	Text   localize.Singular
	Plural localize.Plural
}

func Func(msgID localize.Singular, plural localize.Plural) string {

	var t localize.Singular

	t = `This is an
multiline string`

	t = "Newline remains\n"
	foo.T.Getf("foo test")
	return t
}
