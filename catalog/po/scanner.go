package po

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"unicode"
)

var (
	rePrevMsgContextComments = regexp.MustCompile(`^#\|\s+msgctxt`)  // #| msgctxt
	rePrevMsgIDComments      = regexp.MustCompile(`^#\|\s+msgid`)    // #| msgid
	rePrevStringLineComments = regexp.MustCompile(`^#\|\s+".*"\s*$`) // #| "message"

	reMsgContext   = regexp.MustCompile(`^msgctxt\s+".*"\s*$`)           // msgctxt
	reMsgID        = regexp.MustCompile(`^msgid\s+".*"\s*$`)             // msgid
	reMsgIDPlural  = regexp.MustCompile(`^msgid_plural\s+".*"\s*$`)      // msgid_plural
	reMsgStr       = regexp.MustCompile(`^msgstr\s*".*"\s*$`)            // msgstr
	reMsgStrPlural = regexp.MustCompile(`^msgstr\s*(\[\d+])\s*".*"\s*$`) // msgstr[0]
	reMsgLine      = regexp.MustCompile(`^\s*".*"\s*$`)                  // "message"

	msgCtxPrefix      = "msgctxt "
	msgIDPrefix       = "msgid "
	msgIDPluralPrefix = "msgid_plural "
	msgStrPrefix      = "msgstr"
	emptyMessageStr   = `msgstr ""`
)

type scanner struct {
	lines []string
	pos   int

	lastToken        token
	whitespaceBuffer bytes.Buffer

	ignoreComments bool
}

func newScanner(content string) *scanner {
	content = strings.ReplaceAll(content, "\r", "")
	return &scanner{
		lines: strings.Split(content, "\n"),
		pos:   0,
	}
}

func (s *scanner) scan() (tok token, lit string) {
	line, err := s.read()
	defer func() {
		s.lastToken = tok
	}()

	if err != nil {
		return eof, ""
	}

	// Short path for better performance
	if len(line) == 0 {
		s.unread()
		return s.scanWhitespace()
	}

	line = strings.TrimRightFunc(line, unicode.IsSpace)
	if line[0] == '#' {
		return s.scanComment(line)
	}

	if tokk, l := s.scanMessage(line); tokk != none {
		return tokk, l
	}

	if len(strings.TrimSpace(line)) == 0 {
		s.unread()
		return s.scanWhitespace()
	}

	return failure, line
}

func (s *scanner) scanMessage(line string) (token, string) {

	// We additionally use the functions of the strings package,
	// as these are more performant than just the regular expressions.

	if strings.HasPrefix(line, msgIDPrefix) && reMsgID.MatchString(line) {
		return msgID, line
	}

	if strings.HasPrefix(line, msgStrPrefix) {
		if line == emptyMessageStr || reMsgStr.MatchString(line) {
			return msgStr, line
		} else if reMsgStrPlural.MatchString(line) {
			return msgStrPlural, line
		} else {
			return failure, line
		}
	}

	if strings.HasPrefix(line, msgIDPluralPrefix) && reMsgIDPlural.MatchString(line) {
		return msgIDPlural, line
	}

	if strings.HasPrefix(line, msgCtxPrefix) && reMsgContext.MatchString(line) {
		return msgContext, line
	}

	if strings.HasPrefix(line, `"`) && strings.HasSuffix(line, `"`) && reMsgLine.MatchString(line) {
		// "..."
		switch s.lastToken {
		case msgID, msgIDLine:
			return msgIDLine, line
		case msgIDPlural, msgIDPluralLine:
			return msgIDPluralLine, line
		case msgStr, msgStrLine:
			return msgStrLine, line
		case msgStrPlural, msgStrPluralLine:
			return msgStrPluralLine, line
		case msgContext, msgContextLine:
			return msgContextLine, line
		}
	}

	return none, line
}

func (s *scanner) scanComment(line string) (token, string) {
	lineLen := len(line)

	// comment without content
	if lineLen == 1 {
		if s.ignoreComments {
			return s.scan()
		}
		return commentTranslator, line
	}

	// special comment without content
	if lineLen == 2 {
		switch line[1] {
		case '.', ':', ',', '|':
			return s.scan()
		}
	}

	if s.ignoreComments {
		switch line[1] {
		case ',':
			// Flag comments will always be read
			break
		default:
			return s.scan()
		}
	}

	switch line[1] {
	case '.':
		return commentExtracted, line
	case ':':
		return commentReference, line
	case ',':
		return commentFlags, line
	case '|':
		// #| "..."
		if rePrevMsgContextComments.MatchString(line) {
			return commentPrevContext, line
		} else if rePrevMsgIDComments.MatchString(line) {
			return commentPrevMsgID, line
		} else if rePrevStringLineComments.MatchString(line) {
			switch s.lastToken {
			case commentPrevContext, commentPrevContextLine:
				return commentPrevContextLine, line
			case commentPrevMsgID, commentPrevMsgIDLine:
				return commentPrevMsgIDLine, line
			default:
				return commentPrevUnknown, line
			}
		}

		return failure, line
	default:
		return commentTranslator, line
	}
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *scanner) scanWhitespace() (tok token, lit string) {
	s.whitespaceBuffer.Reset()
	currentLine, _ := s.read()
	s.whitespaceBuffer.WriteString(currentLine)

	for {
		line, err := s.read()
		if err == io.EOF {
			break
		}

		if len(line) == 0 || strings.TrimSpace(line) == "" {
			s.whitespaceBuffer.WriteString(line)
			continue
		}

		s.unread()
		break
	}

	return whitespace, s.whitespaceBuffer.String()
}

func (s *scanner) read() (string, error) {
	if s.pos >= len(s.lines) {
		return "", io.EOF
	}

	line := s.lines[s.pos]
	s.pos++
	return line, nil
}

func (s *scanner) currentLine() string {
	if s.pos > len(s.lines) {
		return ""
	}

	return s.lines[s.pos-1]
}

func (s *scanner) unread() {
	if s.pos > 0 {
		s.pos--
	}
}
