package cldrplural

import (
	"math"
	"strconv"
	"strings"
)

type Category byte

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

type Operand int

const (
	OperandN Operand = iota // the absolute value of N.*
	OperandI                // the integer digits of N.*
	OperandV                // the number of visible fraction digits in N, with trailing zeros.*
	OperandW                // the number of visible fraction digits in N, without trailing zeros.*
	OperandF                // the visible fraction digits in N, with trailing zeros, expressed as an integer.*
	OperandT                // the visible fraction digits in N, without trailing zeros, expressed as an integer.*
	OperandC                // compact decimal exponent value: exponent of the power of 10 used in compact decimal formatting.
)

func (op Operand) String() string {
	if name, ok := OperandNames[op]; ok {
		return name
	}

	return "unknown operand"
}

var OperandNames = map[Operand]string{
	OperandN: "n",
	OperandI: "i",
	OperandV: "v",
	OperandW: "w",
	OperandF: "f",
	OperandT: "t",
	OperandC: "c",
}

var OperandMap = map[string]Operand{
	"n": OperandN,
	"i": OperandI,
	"v": OperandV,
	"w": OperandW,
	"f": OperandF,
	"t": OperandT,
	"c": OperandC,
}

type FormFunc func(ops *Operands) Category

type RuleSet struct {
	Categories []Category
	FormFunc   FormFunc
}

type Operands struct {
	N float64
	I int64
	V int64
	W int64
	F int64
	T int64
	C int64
}

// ExtractOperands converts the representation of a float value into the appropriate operands.
func ExtractOperands(raw string) *Operands {
	op := &Operands{}

	if strings.Contains(raw, "c") {
		cIdx := strings.Index(raw, "c")
		c, _ := strconv.Atoi(raw[cIdx+1:])
		op.C = int64(c)
		raw = formatExponent(raw[:cIdx], c)
	}

	src, _ := strconv.ParseFloat(raw, 64)

	op.N = math.Abs(src)
	op.I = int64(src)

	if pointIdx := strings.Index(raw, "."); pointIdx >= 0 {
		fractionDigits := raw[pointIdx+1:]
		op.V = int64(len(fractionDigits))
		if i, err := strconv.Atoi(fractionDigits); err == nil {
			op.F = int64(i)
		}

		withoutZeros := strings.TrimRight(fractionDigits, "0")
		op.W = int64(len(withoutZeros))
		if i, err := strconv.Atoi(withoutZeros); err == nil {
			op.T = int64(i)
		}
	}

	return op
}

func formatExponent(raw string, c int) string {
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
