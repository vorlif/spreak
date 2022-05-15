package humanize_test

import (
	"fmt"
	"time"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
	"github.com/vorlif/spreak/humanize/locale/ar"
	"github.com/vorlif/spreak/humanize/locale/be"
	"github.com/vorlif/spreak/humanize/locale/de"
	"github.com/vorlif/spreak/humanize/locale/es"
	"github.com/vorlif/spreak/humanize/locale/zhHans"
)

func ExampleHumanizer_Apnumber() {
	parcel := humanize.MustNew(humanize.WithLocale(es.New(), de.New(), zhHans.New()))

	for _, tag := range []language.Tag{language.English, language.German, language.Spanish, language.SimplifiedChinese} {
		h := parcel.CreateHumanizer(tag)
		fmt.Println(h.Apnumber(5))
	}
	// Output:
	// five
	// fünf
	// cinco
	// 五
}

func ExampleHumanizer_Intword() {
	parcel := humanize.MustNew(humanize.WithLocale(es.New(), de.New(), zhHans.New()))

	for _, tag := range []language.Tag{language.English, language.Spanish, language.SimplifiedChinese} {
		h := parcel.CreateHumanizer(tag)
		fmt.Println(h.Intword(1_000_000_000))
	}
	// Output:
	// 1.0 billion
	// 1,0 millardo
	// 1.0 十亿
}

func ExampleHumanizer_Intcomma() {
	parcel := humanize.MustNew(humanize.WithLocale(es.New(), de.New(), zhHans.New()))

	for _, tag := range []language.Tag{language.English, language.Spanish, language.SimplifiedChinese} {
		h := parcel.CreateHumanizer(tag)
		fmt.Println(h.Intcomma(1_000_000_000))
	}
	// Output:
	// 1,000,000,000
	// 1.000.000.000
	// 1,000,000,000
}

func ExampleHumanizer_Ordinal() {
	parcel := humanize.MustNew(humanize.WithLocale(es.New(), de.New(), zhHans.New()))

	for _, tag := range []language.Tag{language.English, language.Spanish, language.SimplifiedChinese} {
		h := parcel.CreateHumanizer(tag)
		fmt.Println(h.Ordinal(5))
	}
	// Output:
	// 5th
	// 5º
	// 第5
}

func ExampleHumanizer_LanguageName() {
	parcel := humanize.MustNew(humanize.WithLocale(es.New(), de.New(), zhHans.New()))

	for _, tag := range []language.Tag{language.English, language.Spanish, language.SimplifiedChinese} {
		h := parcel.CreateHumanizer(tag)
		fmt.Printf("Examples in: %s, %s", h.LanguageName("English"), h.LanguageName("Spanish"))
		fmt.Printf(" and %s\n", h.LanguageName("Simplified Chinese"))
	}
	// Output:
	// Examples in: English, Spanish and Simplified Chinese
	// Examples in: Inglés, Español and Chino simplificado
	// Examples in: 英语, 西班牙语 and 简体中文
}

func ExampleHumanizer_NaturalDay() {
	parcel := humanize.MustNew(humanize.WithLocale(es.New(), de.New(), zhHans.New()))

	for _, tag := range []language.Tag{language.English, language.Spanish, language.SimplifiedChinese} {
		h := parcel.CreateHumanizer(tag)
		fmt.Println(h.NaturalDay(time.Now()))
	}
	// Output:
	// today
	// hoy
	// 今天
}

func ExampleHumanizer_NaturalTime() {
	parcel := humanize.MustNew(humanize.WithLocale(es.New(), ar.New(), zhHans.New()))

	for _, tag := range []language.Tag{language.English, language.Arabic, language.SimplifiedChinese} {
		h := parcel.CreateHumanizer(tag)
		t := time.Now().Add(5 * time.Minute)
		fmt.Println(h.NaturalTime(t))
	}
	// Output:
	// 5 minutes from now
	// ٥ دقائق من الآن
	// 5分钟以后
}

func ExampleHumanizer_NaturalTime_past() {
	parcel := humanize.MustNew(humanize.WithLocale(es.New(), ar.New(), zhHans.New()))

	for _, tag := range []language.Tag{language.English, language.Arabic, language.SimplifiedChinese} {
		h := parcel.CreateHumanizer(tag)
		t := time.Now().Add(-2*time.Hour - 30*time.Minute)
		fmt.Println(h.NaturalTime(t))
	}
	// Output:
	// 2 hours ago
	// منذ ٢ ساعة
	// 2小时之前
}

func ExampleHumanizer_TimeSince() {
	parcel := humanize.MustNew(humanize.WithLocale(es.New(), be.New(), de.New()))

	for _, tag := range []language.Tag{language.English, language.MustParse("be"), language.German} {
		h := parcel.CreateHumanizer(tag)
		t := time.Now().Add(-37 * time.Hour)
		fmt.Println(h.TimeSince(t))
	}
	// Output:
	// 1 day, 13 hours
	// 1 дзень, 13 гадзін
	// 1 Tag, 13 Stunden
}

func ExampleHumanizer_TimeSince_duration() {
	parcel := humanize.MustNew(humanize.WithLocale(es.New(), be.New(), de.New()))

	for _, tag := range []language.Tag{language.English, language.MustParse("be"), language.German} {
		h := parcel.CreateHumanizer(tag)
		t := -80 * time.Hour
		fmt.Println(h.TimeSince(t))
	}
	// Output:
	// 3 days, 8 hours
	// 3  дзён, 8 гадзін
	// 3 Tage, 8 Stunden
}

func ExampleHumanizer_TimeUntil() {
	parcel := humanize.MustNew(humanize.WithLocale(es.New(), be.New(), de.New()))

	for _, tag := range []language.Tag{language.English, language.MustParse("be"), language.German} {
		h := parcel.CreateHumanizer(tag)
		t := time.Now().Add(37 * time.Hour)
		fmt.Println(h.TimeUntil(t))
	}
	// Output:
	// 1 day, 13 hours
	// 1 дзень, 13 гадзін
	// 1 Tag, 13 Stunden
}

func ExampleHumanizer_FormatTime() {
	parcel := humanize.MustNew(humanize.WithLocale(es.New(), be.New(), de.New()))

	now := time.Date(2022, 05, 15, 18, 0, 0, 0, time.Local)
	for _, tag := range []language.Tag{language.English, language.MustParse("be"), language.German} {
		h := parcel.CreateHumanizer(tag)
		fmt.Println(h.FormatTime(now, humanize.DateTimeFormat))
	}
	// Output:
	// May 15, 2022, 6 p.m.
	// траўня 15, 2022, 6 папаўдні
	// 15. Mai 2022 18:00
}

func ExampleHumanizer_FormatTime_shortDate() {
	parcel := humanize.MustNew(humanize.WithLocale(es.New(), zhHans.New(), de.New()))

	now := time.Date(2022, 05, 15, 18, 0, 0, 0, time.Local)
	for _, tag := range []language.Tag{language.English, language.SimplifiedChinese, language.German} {
		h := parcel.CreateHumanizer(tag)
		fmt.Println(h.FormatTime(now, humanize.ShortDateFormat))
	}
	// Output:
	// 05/15/2022
	// 2022年5月15日
	// 15.05.2022
}

func ExampleHumanizer_FormatTime_shortDateTime() {
	parcel := humanize.MustNew(humanize.WithLocale(es.New(), zhHans.New(), de.New()))

	now := time.Date(2022, 05, 15, 18, 0, 0, 0, time.Local)
	for _, tag := range []language.Tag{language.English, language.SimplifiedChinese, language.German} {
		h := parcel.CreateHumanizer(tag)
		fmt.Println(h.FormatTime(now, humanize.ShortDatetimeFormat))
	}
	// Output:
	// 05/15/2022 6 p.m.
	// 2022年5月15日 18:00
	// 15.05.2022 18:00
}
