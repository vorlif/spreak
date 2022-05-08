package plural

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

const (
	prefixNplurals = "nplurals"
	prefixPlural   = "plural"
	variableN      = "n"
)

type scanner struct {
	r *bufio.Reader
}

func newScanner(r io.Reader) *scanner {
	return &scanner{r: bufio.NewReader(r)}
}

func (s *scanner) scan() (tok token, lit string) {
	ch := s.read()
	if ch == scannerEOF {
		return eof, ""
	}

	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		return s.scanText()
	} else if isDigit(ch) {
		s.unread()
		return s.scanNumber()
	}

	switch ch {
	case scannerEOF:
		return eof, ""
	case '%':
		return reminder, string(ch)
	case '?':
		return question, string(ch)
	case ':':
		return colon, string(ch)
	case ';':
		return semicolon, string(ch)
	case '(':
		return leftBracket, string(ch)
	case ')':
		return rightBracket, string(ch)
	case '&', '!', '|', '<', '>', '=':
		nextCh := s.read()
		switch string([]rune{ch, nextCh}) {
		case "!=":
			return notEqual, "!="
		case "&&":
			return logicalAnd, "&&"
		case "==":
			return equal, "=="
		case "||":
			return logicalOr, "||"
		case ">=":
			return greaterOrEqual, ">="
		case "<=":
			return lessOrEqual, "<="
		}

		s.unread()
		switch ch {
		case '<':
			return less, string(ch)
		case '>':
			return greater, string(ch)
		case '=':
			return assign, string(ch)
		}
	}

	return failure, string(ch)
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *scanner) scanWhitespace() (tok token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == scannerEOF {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return whitespace, buf.String()
}

func (s *scanner) scanNumber() (tok token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == scannerEOF {
			break
		} else if !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return number, buf.String()
}

func (s *scanner) scanText() (tok token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == scannerEOF {
			break
		} else if !isLetter(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	switch strings.ToLower(buf.String()) {
	case prefixNplurals:
		return nPlurals, buf.String()
	case prefixPlural:
		return plural, buf.String()
	case variableN:
		return variable, buf.String()
	default:
		return failure, buf.String()
	}
}

func (s *scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return scannerEOF
	}
	return ch
}

func (s *scanner) unread() { _ = s.r.UnreadRune() }

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

func isDigit(ch rune) bool { return ch >= '0' && ch <= '9' }

// eof represents a marker rune for the end of the reader.
var scannerEOF = rune(0)
