package ar

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("ar"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:      "j F، Y",
			TimeFormat:      "g:i A",
			YearMonthFormat: "F Y",
			MonthDayFormat:  "j F",
			ShortDateFormat: "d‏/m‏/Y",
		},
	}
}
