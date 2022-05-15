package fi

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("fi"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          "j. E Y",
			TimeFormat:          "G.i",
			DateTimeFormat:      `j. E Y \k\e\l\l\o G.i`,
			YearMonthFormat:     "F Y",
			MonthDayFormat:      "j. F",
			ShortDateFormat:     "j.n.Y",
			ShortDatetimeFormat: "j.n.Y G.i",
			FirstDayOfWeek:      1,
		},
	}
}
