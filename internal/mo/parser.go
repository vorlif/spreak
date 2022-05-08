package mo

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"strings"

	"github.com/vorlif/spreak/internal/util"
)

const (
	magicLittleEndian = 0x950412de
	magicBigEndian    = 0xde120495

	eotSeparator = "\x04" // msgctxt and msgid separator
	nulSeparator = "\x00" // msgid and msgstr separator
)

type parser struct {
	bo binary.ByteOrder
	r  io.ReadSeeker

	magicNumber  uint32
	majorVersion uint16
	minorVersion uint16
	msgIDCount   uint32
	msgIDOffset  uint32
	msgStrOffset uint32
	hashSize     uint32
	hashOffset   uint32
	messageData  []*msgData
}

type msgData struct {
	idOffset          uint32
	idLength          uint32
	translationOffset uint32
	translationLength uint32
}

func ParseBytes(data []byte) (*File, error) {
	return newParser(bytes.NewReader(data)).parse()
}

func ParseReader(r io.ReadSeeker) (*File, error) {
	return newParser(r).parse()
}

func newParser(r io.ReadSeeker) *parser {
	return &parser{
		r:  r,
		bo: binary.LittleEndian,
	}
}

func (p *parser) parse() (*File, error) {
	if err := p.parseByteOrder(); err != nil {
		return nil, err
	}

	if err := p.parseHeader(); err != nil {
		return nil, err
	}

	if err := p.parseMessageData(); err != nil {
		return nil, err
	}

	file, err := p.parseMessages()
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (p *parser) parseByteOrder() error {
	var magicNumber uint32
	if err := p.binaryRead(&magicNumber); err != nil {
		return err
	}

	switch magicNumber {
	case magicLittleEndian:
		p.bo = binary.LittleEndian
	case magicBigEndian:
		p.bo = binary.BigEndian
	default:
		return errors.New("mo file has invalid magic number")
	}

	p.magicNumber = magicNumber
	return nil
}

func (p *parser) parseHeader() error {
	var header struct {
		MajorVersion uint16
		MinorVersion uint16
		MsgIDCount   uint32
		MsgIDOffset  uint32
		MsgStrOffset uint32
		HashSize     uint32
		HashOffset   uint32
	}
	if err := p.binaryRead(&header); err != nil {
		return err
	}
	if v := header.MajorVersion; v != 0 && v != 1 {
		return errors.New("mo file has invalid major version number")
	}
	if v := header.MinorVersion; v != 0 && v != 1 {
		return errors.New("mo file has invalid minor version number")
	}

	p.majorVersion = header.MajorVersion
	p.minorVersion = header.MinorVersion
	p.msgIDCount = header.MsgIDCount
	p.msgIDOffset = header.MsgIDOffset
	p.msgStrOffset = header.MsgStrOffset
	p.hashSize = header.HashSize
	p.hashOffset = header.HashOffset
	return nil
}

func (p *parser) parseMessageData() error {
	p.messageData = make([]*msgData, p.msgIDCount)

	if _, err := p.r.Seek(int64(p.msgIDOffset), 0); err != nil {
		return err
	}

	for i := 0; i < len(p.messageData); i++ {
		p.messageData[i] = &msgData{}
	}

	for _, d := range p.messageData {
		if err := p.binaryRead(&d.idLength); err != nil {
			return err
		}

		if err := p.binaryRead(&d.idOffset); err != nil {
			return err
		}
	}

	if _, err := p.r.Seek(int64(p.msgStrOffset), 0); err != nil {
		return err
	}

	for _, d := range p.messageData {
		if err := p.binaryRead(&d.translationLength); err != nil {
			return err
		}

		if err := p.binaryRead(&d.translationOffset); err != nil {
			return err
		}
	}

	return nil
}

func (p *parser) parseMessages() (*File, error) {
	file := &File{
		Header:   Header{},
		Messages: make([]Message, 0, len(p.messageData)),
	}

	for _, d := range p.messageData {
		id, localized, err := p.parseMessage(d)
		if err != nil {
			return nil, err
		}

		var msg = Message{
			ID:  id,
			Str: localized,
		}

		// Header
		if id == "" {
			file.Header = p.parseMessageHeader(&msg)
			continue
		}

		// Context
		if idx := strings.Index(msg.ID, eotSeparator); idx != -1 {
			msg.Context, msg.ID = msg.ID[:idx], msg.ID[idx+1:]
		}

		// Plural
		if idx := strings.Index(msg.ID, nulSeparator); idx != -1 {
			msg.ID, msg.IDPlural = msg.ID[:idx], msg.ID[idx+1:]
			msg.StrPlural = strings.Split(msg.Str, nulSeparator)
			msg.Str = ""
		}

		file.Messages = append(file.Messages, msg)
	}

	return file, nil
}

func (p *parser) parseMessage(data *msgData) (id string, localized string, err error) {
	if _, err = p.r.Seek(int64(data.idOffset), 0); err != nil {
		return
	}
	tmp := make([]byte, data.idLength)
	if err = p.read(tmp); err != nil {
		return
	}
	id = string(tmp)

	if _, err = p.r.Seek(int64(data.translationOffset), 0); err != nil {
		return
	}

	tmp = make([]byte, data.translationLength)
	if err = p.read(tmp); err != nil {
		return
	}
	localized = string(tmp)
	return
}

func (p *parser) parseMessageHeader(msg *Message) (header Header) {
	if msg.Str == "" {
		return
	}

	for _, line := range strings.Split(msg.Str, "\n") {
		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		switch strings.ToUpper(key) {
		case strings.ToUpper(util.HeaderProjectIDVersion):
			header.ProjectIDVersion = val
		case strings.ToUpper(util.HeaderReportMsgIDBugsTo):
			header.ReportMsgidBugsTo = val
		case strings.ToUpper(util.HeaderPOTCreationDate):
			header.POTCreationDate = val
		case strings.ToUpper(util.HeaderPORevisionDate):
			header.PORevisionDate = val
		case strings.ToUpper(util.HeaderLastTranslator):
			header.LastTranslator = val
		case strings.ToUpper(util.HeaderLanguageTeam):
			header.LanguageTeam = val
		case strings.ToUpper(util.HeaderLanguage):
			header.Language = val
		case strings.ToUpper(util.HeaderMIMEVersion):
			header.MimeVersion = val
		case strings.ToUpper(util.HeaderContentType):
			header.ContentType = val
		case strings.ToUpper(util.HeaderContentTransferEncoding):
			header.ContentTransferEncoding = val
		case strings.ToUpper(util.HeaderPluralForms):
			header.PluralForms = val
		case strings.ToUpper(util.HeaderXGenerator):
			header.XGenerator = val
		default:
			if header.UnknownFields == nil {
				header.UnknownFields = make(map[string]string)
			}
			header.UnknownFields[key] = val
		}
	}

	return
}

func (p *parser) binaryRead(data interface{}) error {
	return binary.Read(p.r, p.bo, data)
}

func (p *parser) read(d []byte) error {
	_, err := p.r.Read(d)
	return err
}
