package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak"
	"github.com/vorlif/spreak/localize"
)

const (
	// Title of the website
	// TRANSLATORS: Spreak is the name of the application
	Title localize.Singular = "Welcome to the Spreak Tour"
)

// xspreak: template
var page = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<body>

<h1>{{.Title}}</h1>

<p>{{ .T.Getf "Nice to see you %s" .User}}
{{range .Paragraphs}}<p>{{.}}</p>{{end}}
</body>
</html>
`))

func main() {
	bundle, err := spreak.NewBundle(
		spreak.WithSourceLanguage(language.English),
		spreak.WithDefaultDomain("httptempl"),
		spreak.WithDomainPath("httptempl", "../../locale"),
		spreak.WithLanguage(language.German, language.French),
	)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lang := r.FormValue("lang")
		accept := r.Header.Get("Accept-Language")

		localizer := spreak.NewLocalizer(bundle, lang, accept)
		username := r.FormValue("name")
		if username == "" {
			username = "John"
		}

		paragraphs := []localize.Singular{
			// TRANSLATORS: spreak is the name of the application
			localizer.Get("In the next steps we will see how we use spreak"),
			localizer.NGet("You will understand why Florian has a dog", "You will understand why Florian has dogs", 1),
		}

		err := page.Execute(w, map[string]interface{}{
			"Title":      Title,
			"User":       username,
			"Paragraphs": paragraphs,
			"T":          localizer,
		})
		if err != nil {
			panic(err)
		}
	})

	fmt.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
