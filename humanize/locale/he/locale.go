package he

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("he"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          "j בF Y",
			TimeFormat:          "H:i",
			DateTimeFormat:      "j בF Y H:i",
			YearMonthFormat:     "F Y",
			MonthDayFormat:      "j בF",
			ShortDateFormat:     "d/m/Y",
			ShortDatetimeFormat: "d/m/Y H:i",
		},
	}
}
