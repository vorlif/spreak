package humanize

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHumanizer_Apnumber(t *testing.T) {
	t.Run("apnumber i18n", func(t *testing.T) {
		h := createGermanHumanizer(t)

		assert.Equal(t, "eins", h.Apnumber(1))
		assert.Equal(t, "zwei", h.Apnumber(2))
		assert.Equal(t, "-1", h.Apnumber(-1))
		assert.Equal(t, "0", h.Apnumber(0))
		assert.Equal(t, "neun", h.Apnumber(9))
		assert.Equal(t, "10", h.Apnumber(10))
		assert.Equal(t, "%!(string=test)", h.Apnumber("test"))

		tests := []struct {
			input interface{}
			want  string
		}{
			{0, "0"},
			{1, "eins"},
			{2, "zwei"},
			{3, "drei"},
			{4, "vier"},
			{5, "fünf"},
			{6, "sechs"},
			{7, "sieben"},
			{8, "acht"},
			{9, "neun"},
			{10, "10"},
			{"test", "%!(string=test)"},
		}
		for _, tt := range tests {
			t.Run(tt.want, func(t *testing.T) {
				assert.Equal(t, tt.want, h.Apnumber(tt.input))
			})
		}
	})

	t.Run("apnumber source", func(t *testing.T) {
		h := createSourceHumanizer(t)
		tests := []struct {
			input interface{}
			want  string
		}{
			{1, "one"},
			{2, "two"},
			{3, "three"},
			{4, "four"},
			{5, "five"},
			{6, "six"},
			{7, "seven"},
			{8, "eight"},
			{9, "nine"},
			{10, "10"},
			{"test", "%!(string=test)"},
		}
		for _, tt := range tests {
			t.Run(tt.want, func(t *testing.T) {
				assert.Equal(t, tt.want, h.Apnumber(tt.input))
			})
		}
	})
}

func TestHumanizer_Intword(t *testing.T) {
	t.Run("intword i18n", func(t *testing.T) {
		h := createGermanHumanizer(t)
		tests := []struct {
			input string
			want  string
		}{
			{"100", "100"},
			{"1000", "1.000"},
			{"1000000", "1,0 Million"},
			{"1200000", "1,2 Millionen"},
			{"1290000", "1,3 Millionen"},
			{"1000000000", "1,0 Milliarde"},
			{"2000000000", "2,0 Milliarden"},
			{"6000000000000", "6,0 Billionen"},
			{"6000000000000", "6,0 Billionen"},
		}
		length := len(tests)
		for i := 0; i < length; i++ {
			negativeTest := tests[i]
			negativeTest.input = "-" + tests[i].input
			negativeTest.want = "-" + tests[i].want
			tests = append(tests, negativeTest)
		}

		for _, tt := range tests {
			t.Run(tt.want, func(t *testing.T) {
				assert.Equal(t, tt.want, h.Intword(tt.input))
			})
		}

		assert.Equal(t, "%!(string=test)", h.Intword("test"))
		assert.Equal(t, "%!(string=-test)", h.Intword("-test"))
	})

	t.Run("intword source", func(t *testing.T) {
		h := createSourceHumanizer(t)
		tests := []struct {
			input string
			want  string
		}{

			{"100", "100"},
			{"1000", "1,000"},
			{"1000000", "1.0 million"},
			{"1200000", "1.2 million"},
			{"1290000", "1.3 million"},
			{"1000000000", "1.0 billion"},
			{"2000000000", "2.0 billion"},
			{"6000000000000", "6.0 trillion"},
			{"1300000000000000", "1.3 quadrillion"},
			{"3500000000000000000000", "3.5 sextillion"},
			{"8100000000000000000000000000000000", "8.1 decillion"},
			{"1" + strings.Repeat("0", 100), "1.0 googol"},
			{"1" + strings.Repeat("0", 104), "100,000,000,000,000,000,191,567,508,573,466,873,621,595,512,726,519,201,115,280,351,459,937,932,420,398,875,596,123,614,510,818,032,353,280"},
		}
		length := len(tests)
		for i := 0; i < length; i++ {
			negativeTest := tests[i]
			negativeTest.input = "-" + tests[i].input
			negativeTest.want = "-" + tests[i].want
			tests = append(tests, negativeTest)
		}

		for _, tt := range tests {
			t.Run(tt.want, func(t *testing.T) {
				assert.Equal(t, tt.want, h.Intword(tt.input))
			})
		}

		assert.Equal(t, "%!(string=test)", h.Intword("test"))
		assert.Equal(t, "%!(string=-test)", h.Intword("-test"))
	})
}

func TestHumanizer_Intcomma(t *testing.T) {
	t.Run("intcomma i18n", func(t *testing.T) {
		h := createGermanHumanizer(t)
		assert.Equal(t, "1.234.567,123457", h.Intcomma("1234567.1234567"))

		tests := []struct {
			input interface{}
			want  string
		}{
			{100, "100"},
			{1000, "1.000"},
			{10123, "10.123"},
			{10311, "10.311"},
			{1000000, "1.000.000"},
			{1234567.25, "1.234.567,25"},
			{"100", "100"},
			{"1000", "1.000"},
			{"10123", "10.123"},
			{"10311", "10.311"},
			{"1000000", "1.000.000"},
			{"1234567.1234567", "1.234.567,123457"},
			{"test", "%!(string=test)"},
		}

		for _, tt := range tests {
			t.Run(tt.want, func(t *testing.T) {
				assert.Equal(t, tt.want, h.Intcomma(tt.input))
			})
		}
	})

	t.Run("intcomma source", func(t *testing.T) {
		tests := []struct {
			input interface{}
			want  string
		}{
			{100, "100"},
			{1000, "1,000"},
			{10123, "10,123"},
			{10311, "10,311"},
			{1000000, "1,000,000"},
			{1234567.25, "1,234,567.25"},
			{"100", "100"},
			{"1000", "1,000"},
			{"10123", "10,123"},
			{"10311", "10,311"},
			{"1000000", "1,000,000"},
			{"1234567.1234567", "1,234,567.123457"},
			{"test", "%!(string=test)"},
		}

		h := createSourceHumanizer(t)
		for _, tt := range tests {
			t.Run(tt.want, func(t *testing.T) {
				assert.Equal(t, tt.want, h.Intcomma(tt.input))
			})
		}
	})
}

func TestHumanizer_Ordinal(t *testing.T) {
	t.Run("ordinal i18n", func(t *testing.T) {
		tests := []struct {
			input interface{}
			want  string
		}{
			{"1", "1."},
			{"2", "2."},
			{"3", "3."},
			{"4", "4."},
			{"11", "11."},
			{"12", "12."},
			{"13", "13."},
			{"101", "101."},
			{"102", "102."},
			{"103", "103."},
			{"111", "111."},
			{"something else", "%!(string=something else)"},
			{nil, "<nil>"},
		}

		h := createGermanHumanizer(t)
		for _, tt := range tests {
			t.Run(tt.want, func(t *testing.T) {
				assert.Equal(t, tt.want, h.Ordinal(tt.input))
			})
		}
	})

	t.Run("ordinal source", func(t *testing.T) {
		tests := []struct {
			input interface{}
			want  string
		}{
			{"1", "1st"},
			{"2", "2nd"},
			{"3", "3rd"},
			{"4", "4th"},
			{"11", "11th"},
			{"12", "12th"},
			{"13", "13th"},
			{"101", "101st"},
			{"102", "102nd"},
			{"103", "103rd"},
			{"111", "111th"},
			{"something else", "%!(string=something else)"},
			{nil, "<nil>"},
		}

		h := createSourceHumanizer(t)
		for _, tt := range tests {
			t.Run(tt.want, func(t *testing.T) {
				assert.Equal(t, tt.want, h.Ordinal(tt.input))
			})
		}
	})
}
