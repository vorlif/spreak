package main

import (
	"html/template"
	"net/http"
)

const Title = "Welcome to the spreak tour"

// In the following text, the strings are extracted by xstreak because it has the tag "xspreak: template".
// xspreak: template
var notFoundTemplate = template.Must(template.New("").Parse(`
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

	err := notFoundTemplate.Execute(w, map[string]interface{}{
		"Title":      localizer.Get(Title),
		"User":       username,
		"Paragraphs": paragraphs,
		"i18n":       NewI18N(localizer),
	})
	if err != nil {
		panic(err)
	}
}
