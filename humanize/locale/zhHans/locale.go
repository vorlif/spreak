package zhHans

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("zh-Hans"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          "Y年n月j日",
			TimeFormat:          "H:i",
			DateTimeFormat:      "Y年n月j日 H:i",
			YearMonthFormat:     "Y年n月",
			MonthDayFormat:      "m月j日",
			ShortDateFormat:     "Y年n月j日",
			ShortDatetimeFormat: "Y年n月j日 H:i",
			FirstDayOfWeek:      1,
		},
	}
}
