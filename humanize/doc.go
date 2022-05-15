// Package humanize provides a collection of functions to convert Go data structures into a human-readable format.
//
// It was adopted in large part by the Django project and is therefore able to translate into several languages.
// A list of all supported languages can be found in the locale package.
//
// Usage
//
// Create a parcel with the languages you want to use.
// 		parcel := humanize.MustNew(
//		    humanize.WithLocale(es.New(), ar.New(), zhHans.New()),
//		)
//
// Create a humanizer
// 		h := parcel.CreateHumanizer(language.Spanish)
//
// Use it
// 		fmt.Println(h.Intword(1_000_000_000))
// 		// Output: 1,0 millardo
package humanize
