package humanize

import (
	"io/fs"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak"
)

const (
	djangoDomain = "django"
)

type Parcel struct {
	bundle  *spreak.Bundle
	locales map[language.Tag]*FormatData
}

type options struct {
	bundleOptions []spreak.BundleOption
	locales       []*LocaleData
}

type Humanizer struct {
	loc    *spreak.Localizer
	format *FormatData
}

type LocaleData struct {
	Lang   language.Tag
	Fs     fs.FS
	Format *FormatData
}

type FormatData struct {
	DateFormat          string
	TimeFormat          string
	DateTimeFormat      string
	YearMonthFormat     string
	MonthDayFormat      string
	ShortDateFormat     string
	ShortDatetimeFormat string
	FirstDayOfWeek      int
}

var fallbackFormat = &FormatData{
	DateFormat:          "N j, Y",
	TimeFormat:          "P",
	DateTimeFormat:      "N j, Y, P",
	YearMonthFormat:     "F Y",
	MonthDayFormat:      "F j",
	ShortDateFormat:     "m/d/Y",
	ShortDatetimeFormat: "m/d/Y P",
	FirstDayOfWeek:      0,
}

type Option func(opts *options) error

func WithLocale(data ...*LocaleData) Option {
	return func(opts *options) error {
		opts.locales = append(opts.locales, data...)
		return nil
	}
}

func WithBundleOption(opt spreak.BundleOption) Option {
	return func(opts *options) error {
		opts.bundleOptions = append(opts.bundleOptions, opt)
		return nil
	}
}

func New(opts ...Option) (*Parcel, error) {
	o := &options{
		bundleOptions: nil,
		locales:       nil,
	}

	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}

	for _, d := range o.locales {
		if d.Format == nil {
			d.Format = fallbackFormat
		}
		d.Format.SetDefaults()
	}

	loader := newLoader(o.locales)

	parcel := &Parcel{
		locales: make(map[language.Tag]*FormatData, len(o.locales)),
	}
	languages := make([]interface{}, 0, len(loader.locales))
	for tag, data := range loader.locales {
		languages = append(languages, tag)
		if data.Format != nil {
			parcel.locales[tag] = data.Format
		}
	}

	o.bundleOptions = append(o.bundleOptions,
		spreak.WithSourceLanguage(language.English),
		spreak.WithDefaultDomain(djangoDomain),
		spreak.WithDomainLoader(djangoDomain, loader),
		spreak.WithLanguage(languages...),
	)

	bundle, err := spreak.NewBundle(o.bundleOptions...)
	if err != nil {
		return nil, err
	}

	parcel.bundle = bundle
	return parcel, nil
}

func MustNew(opts ...Option) *Parcel {
	parcel, err := New(opts...)
	if err != nil {
		panic(err)
	}
	return parcel
}

func (p *Parcel) CreateHumanizer(lang language.Tag) *Humanizer {
	loc := spreak.NewLocalizer(p.bundle, lang)

	if data, ok := p.locales[loc.Language()]; ok {
		return &Humanizer{loc: loc, format: data}
	}

	return &Humanizer{loc: loc, format: fallbackFormat}
}

func (f *FormatData) SetDefaults() {
	if f.DateFormat == "" {
		f.DateFormat = fallbackFormat.DateFormat
	}
	if f.TimeFormat == "" {
		f.TimeFormat = fallbackFormat.TimeFormat
	}
	if f.DateTimeFormat == "" {
		f.DateTimeFormat = fallbackFormat.DateTimeFormat
	}
	if f.YearMonthFormat == "" {
		f.YearMonthFormat = fallbackFormat.YearMonthFormat
	}
	if f.MonthDayFormat == "" {
		f.MonthDayFormat = fallbackFormat.MonthDayFormat
	}
	if f.ShortDateFormat == "" {
		f.ShortDateFormat = fallbackFormat.ShortDateFormat
	}
	if f.ShortDatetimeFormat == "" {
		f.ShortDatetimeFormat = fallbackFormat.ShortDatetimeFormat
	}
}
