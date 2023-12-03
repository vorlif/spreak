package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak"
)

const CookieName = "lang"

func main() {
	// Create a bundle to load the translations
	bundle, err := spreak.NewBundle(
		spreak.WithFallbackLanguage(language.English),
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

// createLocalizer creates a localizer for a request which can be used to translate texts.
func (h *handler) createLocalizer(r *http.Request) *spreak.KeyLocalizer {
	clientLanguages := make([]any, 0)

	if cookie, err := r.Cookie(CookieName); err == nil && cookie.Value != "" {
		clientLanguages = append(clientLanguages, cookie.Value)
	}

	clientLanguages = append(clientLanguages, r.Header.Get("Accept-Language"))
	return spreak.NewKeyLocalizer(h.bundle, clientLanguages...)
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

	err := h.templates.ExecuteTemplate(w, name, map[string]any{
		"i18n": NewI18N(localizer),
	})
	if err != nil {
		panic(err)
	}
}

// setLanguage sets a cookie for a selected language.
// This cookie is used for future requests to set the language of translations.
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

// We wrap the localizer to use our own function names.
type i18n struct {
	loc *spreak.KeyLocalizer
}

func NewI18N(loc *spreak.KeyLocalizer) *i18n { return &i18n{loc: loc} }

// We want to use Tr instead of "Getf" and therefore wrap the method.
func (in *i18n) Tr(message string, vars ...any) string { return in.loc.Getf(message, vars...) }

func (in *i18n) TrN(message string, a any, vars ...any) string {
	return in.loc.NGetf(message, a, vars...)
}
