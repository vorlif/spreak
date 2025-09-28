package po

import (
	"bufio"
	"io"
	"regexp"
	"strings"
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
	lineScanner *bufio.Scanner
	// lastToken is the last token read
	lastToken token
	// pos is the current line number
	pos int
	n   int

	ignoreComments bool
}

func newScanner(r io.Reader) *scanner {
	// content = bytes.ReplaceAll(content, []byte("\r"), []byte(""))
	return &scanner{
		lineScanner: bufio.NewScanner(r),
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
		return s.scanWhitespace()
	}

	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return s.scanWhitespace()
	}

	if line[0] == '#' {
		return s.scanComment(line)
	}

	if tokk, l := s.scanMessage(line); tokk != none {
		return tokk, l
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
	for {
		line, err := s.read()
		if err == io.EOF {
			break
		}

		if len(line) == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		s.unread()
		break
	}

	return whitespace, ""
}

func (s *scanner) read() (string, error) {
	if s.n == 1 {
		s.n = 0
		return s.lineScanner.Text(), nil
	}

	if !s.lineScanner.Scan() {
		return "", io.EOF
	}
	if s.lineScanner.Err() != nil {
		return "", s.lineScanner.Err()
	}

	s.pos++
	return s.lineScanner.Text(), nil
}

func (s *scanner) currentLine() string { return s.lineScanner.Text() }
func (s *scanner) unread()             { s.n = 1 }
