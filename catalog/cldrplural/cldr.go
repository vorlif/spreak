package cldrplural

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/vorlif/spreak/internal/util"
)

type Category int

const (
	Zero Category = iota
	One
	Two
	Few
	Many
	Other
)

var CategoryNames = map[Category]string{
	Zero:  "Zero",
	One:   "One",
	Two:   "Two",
	Few:   "Few",
	Many:  "Many",
	Other: "Other",
}

func (cat Category) String() string {
	if name, ok := CategoryNames[cat]; ok {
		return name
	}

	return "unknown"
}

type FormFunc func(ops *Operands) Category

// RuleSet Represents a collection of plural rules for a language.
type RuleSet struct {
	// The categories that can be returned by for these rules.
	Categories []Category
	// A function that returns the associated category for a value.
	FormFunc FormFunc
}

func (rs *RuleSet) Evaluate(a any) (Category, error) {
	ops, err := NewOperands(a)
	if err != nil {
		return Other, err
	}

	return rs.FormFunc(ops), nil
}

// The Operands are numeric values corresponding to features of the source number.
//
// See: http://unicode.org/reports/tr35/tr35-numbers.html#Plural_Operand_Meanings
type Operands struct {
	// The absolute value of the source number.
	N float64
	// The integer digits of the source number.
	I int64
	// The number of visible fraction digits in the source number, with trailing zeros.
	V int64
	// The number of visible fraction digits in the source number, without trailing zeros.
	W int64
	// The visible fraction digits in the source, with trailing zeros, expressed as an integer.*
	F int64
	// The visible fraction digits in source number, without trailing zeros, expressed as an integer.*
	T int64
	// Compact decimal exponent value: exponent of the power of 10 used in compact decimal formatting.
	C int64
}

// MustNewOperands is like NewOperands, but panics if the given value cannot be parsed.
// It simplifies safe initialization of Operands values.
func MustNewOperands(a any) *Operands {
	ops, err := NewOperands(a)
	if err != nil {
		panic(err)
	}
	return ops
}

// NewOperands converts the representation of a float value into the appropriate Operands.
func NewOperands(a any) (*Operands, error) {
	a = util.Indirect(a)
	if a == nil {
		return nil, errors.New("operands value is nil")
	}

	switch v := a.(type) {
	case string:
		return newOperandsString(v)
	case int64:
		return newOperandsInt(v), nil
	case int:
		return newOperandsInt(int64(v)), nil
	case float32:
		return newOperandsString(fmt.Sprintf("%v", v))
	case float64:
		return newOperandsString(fmt.Sprintf("%v", v))
	default:
		num, err := util.ToNumber(v)
		if err != nil {
			return nil, err
		}
		return newOperandsString(fmt.Sprintf("%v", num))
	}
}

func newOperandsInt(i int64) *Operands {
	if i < 0 {
		i = -i
	}
	return &Operands{float64(i), i, 0, 0, 0, 0, 0}
}

func newOperandsString(raw string) (*Operands, error) {
	op := &Operands{}

	if cIdx := strings.Index(raw, "c"); cIdx >= 0 {
		c, err := strconv.Atoi(raw[cIdx+1:])
		if err != nil {
			return nil, err
		}
		op.C = int64(c)
		raw = shiftDecimalPoint(raw[:cIdx], c)
	}

	src, errP := strconv.ParseFloat(raw, 64)
	if errP != nil {
		return nil, errP
	}

	op.N = math.Abs(src)
	op.I = int64(src)

	if pointIdx := strings.Index(raw, "."); pointIdx >= 0 {
		fractionDigits := raw[pointIdx+1:]
		if fractionDigits != "" {
			op.V = int64(len(fractionDigits))
			f, err := strconv.ParseInt(fractionDigits, 10, 64)
			if err != nil {
				return nil, err
			}
			op.F = f
		}

		withoutZeros := strings.TrimRight(fractionDigits, "0")
		if withoutZeros != "" {
			op.W = int64(len(withoutZeros))
			t, err := strconv.ParseInt(withoutZeros, 10, 64)
			if err != nil {
				return nil, err
			}

			op.T = t
		}
	}

	return op, nil
}

func shiftDecimalPoint(raw string, c int) string {
	var s strings.Builder

	shift := false
	for _, r := range raw {
		if r == '.' {
			shift = true
			continue
		}
		if c == 0 && shift {
			s.WriteRune('.')
			shift = false
		}
		if shift {
			c--
		}
		s.WriteRune(r)
	}

	for i := 0; i < c; i++ {
		s.WriteRune('0')
	}
	return s.String()
}
