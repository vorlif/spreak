package ml

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("ml"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          "N j, Y",
			TimeFormat:          "P",
			DateTimeFormat:      "N j, Y, P",
			YearMonthFormat:     "F Y",
			MonthDayFormat:      "F j",
			ShortDateFormat:     "m/d/Y",
			ShortDatetimeFormat: "m/d/Y P",
			FirstDayOfWeek:      0,
		},
	}
}
