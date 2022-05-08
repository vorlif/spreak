// Package spreak provides a simple translation facility based on the concepts of gettext.
//
// Create a bundle with the translations for your domains and your desired languages.
//     bundle, err := spreak.NewBundle(
//         spreak.WithDomainPath(spreak.NoDomain, "../locale"),
//         spreak.WithLanguage(language.German, language.Spanish, language.Chinese),
//     )
//
// Create a Locale or a Localizer
//     t := spreak.NewLocalizer(bundle, language.Spanish)
//
// Use the Localizer to translate messages.
//     fmt.Println(t.Get("Hello world"))
//
// Fundamentals
//
// Domain: A message domain is a set of translatable msgid messages.
// Usually, every software package has its own message domain.
// The domain name is used to determine the message catalog where the translation is looked up.
//
// Default domain: The default domain is used if a domain is not explicitly specified for a requested translation.
// If no default domain is specified, the default domain of the bundle is used.
// If this was not specified either, the domain is NoDomain (an empty string).
//
// Context: Context can be added to strings to be translated.
// A context dependent translation lookup is when a translation for a given string is searched,
// that is limited to a given context. The translation for the same string in a different context can be different.
// The different translations of the same string in different contexts can be stored in the in the same MO file,
// and can be edited by the translator in the same PO file.
// The Context string is visible in the PO file to the translator.
// You should try to make it somehow canonical and never changing.
// Because every time you change an Context, the translator will have to review the translation of msgid.
package spreak
