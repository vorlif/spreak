package hr

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("hr"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          "j. E Y.",
			TimeFormat:          "H:i",
			DateTimeFormat:      "j. E Y. H:i",
			YearMonthFormat:     "F Y.",
			MonthDayFormat:      "j. F",
			ShortDateFormat:     "j.m.Y.",
			ShortDatetimeFormat: "j.m.Y. H:i",
			FirstDayOfWeek:      1,
		},
	}
}
