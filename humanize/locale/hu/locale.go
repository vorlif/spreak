package hu

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("hu"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          "Y. F j.",
			TimeFormat:          "H:i",
			DateTimeFormat:      "Y. F j. H:i",
			YearMonthFormat:     "Y. F",
			MonthDayFormat:      "F j.",
			ShortDateFormat:     "Y.m.d.",
			ShortDatetimeFormat: "Y.m.d. H:i",
			FirstDayOfWeek:      1,
		},
	}
}
