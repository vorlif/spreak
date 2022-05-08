package main

import (
	"golang.org/x/text/language"

	"github.com/vorlif/testdata/sub"

	"github.com/vorlif/spreak"
	sp "github.com/vorlif/spreak"
	alias "github.com/vorlif/spreak/localize"
)

func noop(sing alias.MsgID, plural alias.Plural, context alias.Context, domain alias.Domain) {

}

func outerFuncDef() {
	f := func(msgid alias.Singular, plural alias.Plural, context alias.Context, domain alias.Domain) {

	}

	// not extracted
	f("f-msgid", "f-plural", "f-context", "f-domain")

	// extracted
	noop("noop-msgid", "noop-plural", "noop-context", "noop-domain")
	sub.Func("submsgid", "subplural")
}

// TRANSLATORS: this is not extracted
func localizerCall(loc *sp.Localizer) {
	// TRANSLATORS: this is extracted
	loc.Getf("localizer func call")
}

func builtInFunctions() {
	bundle, err := spreak.NewBundle(
		spreak.WithDefaultDomain(spreak.NoDomain),
		spreak.WithDomainPath(spreak.NoDomain, "./"),
		spreak.WithLanguage(language.English),
	)
	if err != nil {
		panic(err)
	}

	t := spreak.NewLocalizer(bundle, "en")
	// TRANSLATORS: Test
	// multiline
	t.Getf("msgid")
	t.NGetf("msgid-n", "pluralid-n", 10, 10)
	t.DGetf("domain-d", "msgid-d")
	t.DNGetf("domain-dn", "msgid-dn", "pluralid-dn", 10)
	t.PGetf("context-pg", "msgid-pg")
	t.NPGetf("context-np", "msgid-np", "pluralid-np", 10)
	t.DPGetf("domain-dp", "context-dp", "singular-dp")
	t.DNPGetf("domain-dnp", "context-dnp", "msgid-dnp", "pluralid-dnp", 10)
}
