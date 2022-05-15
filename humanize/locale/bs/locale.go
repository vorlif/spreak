package bs

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("bs"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:      "j. N Y.",
			TimeFormat:      "G:i",
			DateTimeFormat:  "j. N. Y. G:i T",
			YearMonthFormat: "F Y.",
			MonthDayFormat:  "j. F",
			ShortDateFormat: "Y M j",
		},
	}
}
