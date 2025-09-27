package po

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type Parser struct {
	s           *scanner
	lastToken   token  // last read token
	lastLiteral string // last read literal
	n           int

	header   *Header
	messages Messages // ctx -> msgID -> Message

	ignoreComments bool
}

func NewParser() *Parser {
	p := &Parser{}

	return p
}

func MustParse(content []byte) *File {
	f, err := Parse(content)
	if err != nil {
		panic(err)
	}

	return f
}

func Parse(content []byte) (*File, error) {
	return ParseString(string(content))
}

func ParseString(content string) (*File, error) {
	p := NewParser()
	return p.Parse(content)
}

// SetIgnoreComments instructs the Parser to skip comments when reading.
// This can lead to better performance when reading Po files if they contain many comments.
// Flags are also read when ignore is activated. Is deactivated by default.
func (p *Parser) SetIgnoreComments(ignore bool) { p.ignoreComments = ignore }

func (p *Parser) reset() {
	p.messages = make(Messages, avgContextsPerFile)
	p.header = nil
}

func (p *Parser) Parse(content string) (*File, error) {
	p.reset()
	p.s = newScanner(content)
	p.s.ignoreComments = p.ignoreComments

	// special case, empty file
	if tok, _ := p.scan(); tok == eof {
		return nil, errors.New("po file cannot be empty")
	}
	p.unscan()

	return p.parseFile()
}

func (p *Parser) parseFile() (*File, error) {

	for {
		msg, errMsg := p.parseMessage()
		if errMsg != nil {
			if errors.Is(errMsg, io.EOF) {
				break
			}
			return nil, errMsg
		}

		if msg.ID == "" {
			if len(msg.Str) == 0 {
				// ignore empty msg
				continue
			}

			if p.header != nil {
				return nil, fmt.Errorf("po file can have only one message with 'msgid \"\"': line %d %q", p.s.pos, p.s.currentLine())
			}

			if header, err := p.parseHeader(msg); err == nil {
				p.header = header
			}

			// ignore header parse err
			continue
		}

		// save message
		if _, ok := p.messages[msg.Context]; !ok {
			exceptedMessages := exceptedMessagesForContext(msg.Context)
			p.messages[msg.Context] = make(map[string]*Message, exceptedMessages)
		}

		if _, ok := p.messages[msg.Context][msg.ID]; ok {
			return nil, fmt.Errorf("po file has a duplicate entry for %q line %d", msg.ID, p.s.pos)
		}

		p.messages[msg.Context][msg.ID] = msg

		if p.lastToken == eof {
			break
		}
	}

	if p.header == nil && len(p.messages) == 0 {
		return nil, errors.New("po file cannot be empty")
	}

	file := &File{Header: &Header{}, Messages: p.messages}
	if p.header != nil {
		file.Header = p.header
	}
	return file, nil
}

func (p *Parser) parseMessage() (*Message, error) {
	if tok, _ := p.scan(); tok == eof {
		return nil, io.EOF
	}

	p.unscan()

	msg := NewMessage()
loop:
	for {
		tok, _ := p.scan()
		switch tok {
		case commentTranslator, commentExtracted, commentReference, commentFlags,
			commentPrevContext, commentPrevMsgID, commentPrevMsgIDLine, commentPrevContextLine:
			p.unscan()
			comment, err := p.parseComment()
			if err != nil {
				return nil, err
			}

			msg.Comment = comment
		default:
			p.unscan()
			break loop
		}
	}

	for {
		tok, lit := p.scan()
		switch tok {
		case failure:
			return nil, p.buildPosError()
		case eof, whitespace, commentTranslator, commentExtracted, commentReference, commentFlags,
			commentPrevContext, commentPrevMsgID, commentPrevMsgIDLine, commentPrevContextLine:
			// Message completed
			return msg, nil
		case msgContext:
			line := lit[8:] // msgctxt
			msg.Context = DecodePoString(line)
		case msgContextLine:
			msg.Context += DecodePoString(lit)
		case msgID:
			line := lit[6:] // msgid
			msg.ID = DecodePoString(line)
		case msgIDLine:
			msg.ID += DecodePoString(lit)
		case msgIDPlural:
			line := lit[13:] // msgid_plural
			msg.IDPlural = DecodePoString(line)
		case msgIDPluralLine:
			msg.IDPlural += DecodePoString(lit)
		case msgStr:
			line := lit[7:] // msgstr
			msg.Str[0] = DecodePoString(line)
		case msgStrLine:
			msg.Str[0] += DecodePoString(lit)
		case msgStrPlural:
			left := strings.Index(lit, "[")
			right := strings.Index(lit, "]")
			// We do not need to check if there is a negative index, because the scanner already checks it.
			idx, err := strconv.Atoi(lit[left+1 : right])
			if err != nil {
				return nil, fmt.Errorf("po file contains an invalid index for a plural translation (line %d)", p.s.pos)
			}

			line := lit[10:] // msgstr[x]
			msg.Str[idx] = DecodePoString(line)
		case msgStrPluralLine:
			lastIdx := len(msg.Str) - 1
			if lastIdx < 0 {
				return nil, p.buildPosError()
			}
			msg.Str[lastIdx] += DecodePoString(lit)
		default:
			return nil, p.buildPosError()
		}
	}
}

func (p *Parser) parseComment() (*Comment, error) {
	comment := NewComment()
	for {
		tok, line := p.scan()

		switch tok {
		case eof:
			return nil, io.EOF
		case failure:
			return nil, p.buildPosError()
		case commentTranslator:
			line = strings.TrimSpace(line[1:]) // #
			if comment.Translator != "" {
				comment.Translator += "\n"
			}
			comment.Translator += line
		case commentExtracted:
			line = strings.TrimSpace(line[2:]) // #.
			if comment.Extracted != "" {
				comment.Extracted += "\n"
			}
			comment.Extracted += line
		case commentFlags:
			line = strings.TrimSpace(line[2:]) // #,
			rawFlags := strings.FieldsFunc(line, func(r rune) bool { return unicode.IsSpace(r) || r == ',' })
			comment.Flags = append(comment.Flags, rawFlags...)
		case commentPrevContext:
			line = strings.TrimSpace(line[10:]) // #| msgctxt
			if line == "" {
				continue
			}
			comment.PrevMsgContext += DecodePoString(line)
		case commentPrevContextLine:
			line = strings.TrimSpace(line[2:]) // #| "..."
			comment.PrevMsgContext += DecodePoString(line)
		case commentPrevMsgID:
			line = strings.TrimSpace(line[8:]) // #| msgid
			if line == "" {
				continue
			}
			comment.PrevMsgID += DecodePoString(line)
		case commentPrevMsgIDLine:
			line = strings.TrimSpace(line[2:]) // #| "..."
			comment.PrevMsgID += DecodePoString(line)
		case commentReference:
			line = strings.TrimSpace(line[2:]) // #:

			rawReferences := strings.Split(line, " ")
			for _, rawRef := range rawReferences {
				rawRef = strings.TrimSpace(rawRef)
				colonIdx := strings.Index(rawRef, ":")
				if colonIdx <= 0 {
					// no line number
					ref := &Reference{
						Path: rawRef,
						Line: -1,
					}
					comment.References = append(comment.References, ref)
					continue
				}

				lineNumber, err := strconv.Atoi(rawRef[colonIdx+1:])
				if err != nil {
					lineNumber = -1
				}
				ref := &Reference{
					Path: rawRef[:colonIdx],
					Line: lineNumber,
				}
				comment.References = append(comment.References, ref)
			}
		default:
			p.unscan()
			return comment, nil
		}
	}
}

func (p *Parser) parseHeader(msg *Message) (*Header, error) {
	header := &Header{}
	lines := strings.Split(msg.Str[0], "\n")
	for _, line := range lines {
		colonIdx := strings.Index(line, ":")
		if colonIdx < 0 {
			if strings.TrimSpace(line) == "" {
				continue
			}

			return nil, fmt.Errorf("po file has an invalid header: %s", lines)
		}

		key := strings.TrimSpace(line[:colonIdx])
		val := strings.TrimSpace(line[colonIdx+1:])
		header.Set(key, val)
	}
	header.Comment = msg.Comment

	return header, nil
}

func (p *Parser) buildPosError() error {
	return fmt.Errorf("po file could not be parsed: line %d %q", p.s.pos, p.s.currentLine())
}

func (p *Parser) scan() (tok token, lit string) {
	if p.n == 1 {
		p.n = 0
		return p.lastToken, p.lastLiteral
	}

	tok, lit = p.s.scan()

	p.lastToken, p.lastLiteral = tok, lit

	return
}

func (p *Parser) unscan() { p.n = 1 }
