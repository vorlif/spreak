package main

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/text/language"

	"github.com/Xuanwo/go-locale"

	"github.com/vorlif/spreak/v2"
	"github.com/vorlif/spreak/v2/localize"
)

const (
	HelloWorldDomain = "helloworld"
	DayInfoDomain    = "dayinfo"
)

var (
	// The texts of Errors can be extracted by xspreak.
	// For this append the -e optiona.
	ErrInvalidHolidayName = errors.New("the name of the holiday has an invalid format")
)

var T *spreak.Localizer

func init() {

	bundle, err := spreak.NewBundle(
		// Specify which language was used in the program code. (Source)
		spreak.WithSourceLanguage(language.English),
		// Sets the domain that will be used if no domain is explicitly specified.
		spreak.WithDefaultDomain(DayInfoDomain),
		// Specifies where the files for the specified domain are stored
		spreak.WithDomainPath(DayInfoDomain, "../locale"),
		spreak.WithDomainPath(HelloWorldDomain, "../locale"),
		// Specifies which languages are to be loaded.
		// If the languages are not found, no error message is returned
		spreak.WithLanguage(language.Spanish, language.French, language.Japanese, language.German),
	)
	if err != nil {
		panic(err)
	}

	// We find out the system language via a library
	if systemLang, err := locale.Detect(); err != nil {
		fmt.Println("System language could not be determined")
		// Fallback - use source language
		T = spreak.NewLocalizer(bundle)
	} else {
		T = spreak.NewLocalizer(bundle, systemLang)
		fmt.Println(T.Getf("System language detected: %v", systemLang))
	}

	// This translation comes from the 'helloworld' domain and should not be extracted.
	// xspreak: ignore
	fmt.Println(T.DGet(HelloWorldDomain, "Hello world"))
}

// These strings are automatically extracted by xspreak.
var weekdays = []localize.MsgID{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

func main() {
	now := time.Now()

	weekdayName := weekdays[now.Weekday()]
	// This string is also automatically extracted by xspreak
	fmt.Printf(T.Get("The weekday is %s\n"), T.Get(weekdayName))

	daysToFriday := (7 + (time.Friday - now.Weekday())) % 7
	// Plural can be output with 'NGet'.
	// xspreak: range: 1..7
	fmt.Println(T.NGetf("%d day left until Friday", "%d days left until Friday", int(daysToFriday), daysToFriday))

	// "Christmas" is automatically extracted because the parameter name is of type localize.MsgID
	christmasHoliday := NewHoliday("Christmas")
	fmt.Printf(T.Get("An example of a holiday is %s\n"), T.Get(christmasHoliday.Name))

	// The empty string is not extracted, instead it is ignored.
	invalidHoliday := NewHoliday("")
	if err := invalidHoliday.IsValid(); err != nil {
		localizedErr := T.LocalizeError(err)
		if !errors.Is(localizedErr, ErrInvalidHolidayName) {
			panic(err)
		}
		fmt.Println(localizedErr)
	}

	// LANGUAGE=de go run main.go

	// OUTPUT:
	// Systemsprache erkannt: de
	// Hallo Welt
	// Der Wochentag ist Montag
	// noch 2 Tage bis Freitag
	// Ein Beispiel für einen Feiertag ist Weihnachten
	// der Name des Feiertags hat ein ungültiges Format
}

func NewHoliday(name localize.MsgID) *Holiday {
	return &Holiday{Name: name}
}

type Holiday struct {
	Name string
}

func (h *Holiday) IsValid() error {
	if h.Name == "" {
		return ErrInvalidHolidayName
	}

	return nil
}
