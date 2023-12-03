package ast

import (
	"io"
	"regexp"
	"strings"
	"unicode"
)

type Token int

const (
	eof Token = iota
	unknown
	Operand     // n, i, v, w, t, f
	Value       // \d+
	Equal       // =
	NotEqual    // !=
	ValueRange  // value..value
	RangeList   // (range | value) (',' range_list)*
	And         // 'and'
	Or          // 'or'
	Remainder   // %
	sample      // '@integer', '@decimal'
	sampleValue // 1, 2, 100.50, etc.
	sampleRange // sampleValue ('~' sampleValue)?
)

var (
	reValue       = regexp.MustCompile(`^\d+$`)
	reRange       = regexp.MustCompile(`^\d+\.\.[1-9]\d*$`)
	reSampleValue = regexp.MustCompile(`^((?:0|[1-9]\d*c?\d*)(?:\.\d*)?|(?:0|[1-9]\d*\d*)(?:\.\d*c?\d*)?)$`)
	reSampleRange = regexp.MustCompile(`^(?:0|[1-9]\d*)(?:\.\d*)?~(?:0|[1-9]\d*)(?:\.\d*)?$`)
)

type scanner struct {
	// The tokens of the rule to be scanned
	tokens []string
	// pos is position of the token at which the scanner is currently located.
	pos int
}

func newScanner(content string) *scanner {
	return &scanner{
		tokens: generateTokens(content),
		pos:    0,
	}
}

func (s *scanner) Scan() (tok Token, lit string) {
	lit, errR := s.read()
	if errR != nil {
		return eof, ""
	}

	switch strings.ToLower(lit) {
	case "":
		return eof, ""
	case "n", "i", "v", "w", "f", "t", "c", "e":
		return Operand, lit
	case "=":
		return Equal, lit
	case "%", "mod":
		return Remainder, "%"
	case "!=":
		return NotEqual, "!="
	case "@integer", "@decimal":
		return sample, lit
	case "and":
		return And, lit
	case "or":
		return Or, lit
	case "…", "...", ".":
		return s.Scan()
	}

	if nextLit, err := s.read(); nextLit == ".." || nextLit == "~" {
		lit += nextLit
		nextLit, err = s.read()
		if err != nil {
			return unknown, lit
		}
		lit += nextLit
	} else if err == nil {
		s.unread()
	}

	if reValue.MatchString(lit) {
		return Value, lit
	} else if reSampleValue.MatchString(lit) {
		return sampleValue, lit
	} else if reRange.MatchString(lit) {
		return ValueRange, lit
	} else if reSampleRange.MatchString(lit) {
		return sampleRange, lit
	}

	return unknown, lit
}

func (s *scanner) read() (string, error) {
	if s.pos >= len(s.tokens) {
		return "", io.EOF
	}

	tok := s.tokens[s.pos]
	s.pos++
	return tok, nil
}

func (s *scanner) unread() {
	if s.pos > 0 {
		s.pos--
	}
}

func generateTokens(content string) []string {
	return strings.FieldsFunc(content, func(r rune) bool {
		return unicode.IsSpace(r) || r == ','
	})
}
