package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak"
	"github.com/vorlif/spreak/localize"
)

const (
	// Title of the website
	// TRANSLATORS: Spreak is the name of the application
	Title localize.Singular = "Welcome to the Spreak Tour"

	CookieName = "lang"
)

func main() {
	// Create a bundle to load the translations
	bundle, err := spreak.NewBundle(
		spreak.WithSourceLanguage(language.English),
		spreak.WithDomainPath(spreak.NoDomain, "locale"),
		spreak.WithLanguage(language.German, language.French, language.Japanese, language.Spanish),
	)
	if err != nil {
		panic(err)
	}

	// Parse the templates
	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}

	// Create a handler for simple routing
	h := &handler{
		bundle:    bundle,
		templates: templates,
	}

	// Start the webserver
	fmt.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", h))
}

// A simple handler which holds the bundle with the translations and the templates
type handler struct {
	bundle    *spreak.Bundle
	templates *template.Template
}

var _ http.Handler = (*handler)(nil)

func (h *handler) createLocalizer(r *http.Request) *spreak.Localizer {
	supportedLangs := make([]interface{}, 0)
	if cookie, err := r.Cookie(CookieName); err == nil && cookie.Value != "" {
		supportedLangs = append(supportedLangs, cookie.Value)
	}
	supportedLangs = append(supportedLangs, r.Header.Get("Accept-Language"))
	return spreak.NewLocalizer(h.bundle, supportedLangs...)
}

// Routing of the requests
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		h.renderTemplate("home.html", w, r)
	case "/forms":
		h.renderTemplate("forms.html", w, r)
	case "/setlang":
		h.setLanguage(w, r)
	default:
		h.NotFound(w, r)
	}
}

// Simple function to render a template.
func (h *handler) renderTemplate(name string, w http.ResponseWriter, r *http.Request) {
	localizer := h.createLocalizer(r)
	err := h.templates.ExecuteTemplate(w, name, map[string]interface{}{"i18n": NewI18N(localizer)})
	if err != nil {
		panic(err)
	}
}

func (h *handler) setLanguage(w http.ResponseWriter, r *http.Request) {
	selectedLang := r.URL.Query().Get("lang")
	if selectedLang != "" {
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: CookieName, Value: selectedLang, Expires: expiration}
		http.SetCookie(w, &cookie)
	} else {
		cookie := http.Cookie{Name: CookieName, Value: "", MaxAge: -1}
		http.SetCookie(w, &cookie)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// In the following text, the strings are extracted by xstreak because it has the tag "xspreak: template".
// xspreak: template
var notFoundTemplates = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html class="h-100">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.2.0-beta1/dist/css/bootstrap.min.css" rel="stylesheet"
          integrity="sha384-0evHe/X+R7YkIZDRvuzKMRqM+OrBnVFBL6DOitfPri4tjfHxaWutUpFmBp4vmVor" crossorigin="anonymous">
</head>
<body>

<main>
    <div class="px-4 py-5 my-5 text-center">
        <img class="d-block mx-auto mb-4" src="https://getbootstrap.com/docs/5.2/assets/brand/bootstrap-logo.svg" alt=""
             width="72" height="57">
        <h1 class="display-5 fw-bold">{{ .i18n.Tr "The page you are looking for does not exist." }}</h1>
        <div class="col-lg-6 mx-auto">
			{{range .Paragraphs}}
				<p class="lead mb-4">{{.}}</p>
			{{end}}
           	<p>{{ .i18n.Tr "Nice to see you %s" .User}}</p>
			<p>{{.Title}}</p>
        </div>
	</div>
</main>

</body>
</html>
`))

// NotFound is a simple handler for not found pages, which gives a nonsense feedback.
func (h *handler) NotFound(w http.ResponseWriter, r *http.Request) {
	localizer := h.createLocalizer(r)
	username := r.FormValue("name")
	if username == "" {
		username = "John"
	}

	paragraphs := []string{
		// TRANSLATORS: spreak is the name of the application
		localizer.Get("In the next steps we will see how we use spreak"),
		localizer.NGet("You will understand why Florian has a dog", "You will understand why Florian has dogs", 1),
	}

	err := notFoundTemplates.Execute(w, map[string]interface{}{
		"Title":      localizer.Get(Title),
		"User":       username,
		"Paragraphs": paragraphs,
		"i18n":       NewI18N(localizer),
	})
	if err != nil {
		panic(err)
	}
}

// We wrap the localizer to use our own function names.
type i18n struct {
	loc *spreak.Localizer
}

func NewI18N(loc *spreak.Localizer) *i18n { return &i18n{loc: loc} }

// We want to use Tr instead of "Getf" and therefore wrap the method.
func (in *i18n) Tr(message string, vars ...interface{}) string { return in.loc.Getf(message, vars...) }
