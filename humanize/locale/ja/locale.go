package ja

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("ja"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          "Y年n月j日",
			TimeFormat:          "G:i",
			DateTimeFormat:      "Y年n月j日G:i",
			YearMonthFormat:     "Y年n月",
			MonthDayFormat:      "n月j日",
			ShortDateFormat:     "Y/m/d",
			ShortDatetimeFormat: "Y/m/d G:i",
		},
	}
}
