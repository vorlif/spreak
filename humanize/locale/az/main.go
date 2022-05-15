package az

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("az"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          "j E Y",
			TimeFormat:          "G:i",
			DateTimeFormat:      "j E Y, G:i",
			YearMonthFormat:     "F Y",
			MonthDayFormat:      "j F",
			ShortDateFormat:     "d.m.Y",
			ShortDatetimeFormat: "d.m.Y H:i",
			FirstDayOfWeek:      1,
		},
	}
}
