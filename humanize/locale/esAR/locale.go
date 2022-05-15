package esAR

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("es-AR"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          `j N Y`,
			TimeFormat:          `H:i`,
			DateTimeFormat:      `j N Y H:i`,
			YearMonthFormat:     `F Y`,
			MonthDayFormat:      `j \d\e F`,
			ShortDateFormat:     `d/m/Y`,
			ShortDatetimeFormat: `d/m/Y H:i`,
			FirstDayOfWeek:      0,
		},
	}
}
