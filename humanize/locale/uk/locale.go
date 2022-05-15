package uk

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("uk"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          "d E Y р.",
			TimeFormat:          "H:i",
			DateTimeFormat:      "d E Y р. H:i",
			YearMonthFormat:     "F Y",
			MonthDayFormat:      "d F",
			ShortDateFormat:     "d.m.Y",
			ShortDatetimeFormat: "d.m.Y H:i",
			FirstDayOfWeek:      1,
		},
	}
}
