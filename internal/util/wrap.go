package util

import (
	"bytes"
	"unicode"
)

const nbsp = 0xA0

func WrapString(s string, pageWidth int) []string {
	lines := make([]string, 0, 10)

	currentLine := &stringBuffer{}
	var currentWordBuf stringBuffer
	var lastSpaceBuf stringBuffer

	for _, char := range s {
		if char == '\n' {
			if currentWordBuf.Len() == 0 {
				lastSpaceBuf.Reset()
			}
			lastSpaceBuf.WriteInto(currentLine)
			currentWordBuf.WriteInto(currentLine)
			lines = append(lines, currentLine.String())
			currentLine.Reset()
		} else if unicode.IsSpace(char) && char != nbsp {
			if currentWordBuf.Len() > 0 { // We had a word before and now a space
				lastSpaceBuf.WriteInto(currentLine)
				currentWordBuf.WriteInto(currentLine)
			}

			lastSpaceBuf.WriteRune(char)
		} else {
			if currentLine.Len()+lastSpaceBuf.Len()+currentWordBuf.Len() >= pageWidth {
				if currentLine.Len() > 0 {
					lines = append(lines, currentLine.String())
					currentLine.Reset()
					lastSpaceBuf.Reset()
				}

				if lastSpaceBuf.Len() > 0 {
					lastSpaceBuf.WriteInto(currentLine)
					currentWordBuf.WriteInto(currentLine)
				}
			}

			currentWordBuf.WriteRune(char)
		}
	}

	if currentWordBuf.Len() > 0 {
		lastSpaceBuf.WriteInto(currentLine)
		currentWordBuf.WriteInto(currentLine)
	}

	lines = append(lines, currentLine.String())
	return lines
}

type stringBuffer struct {
	buff bytes.Buffer
	len  int
}

func (b *stringBuffer) Len() int {
	return b.len
}

func (b *stringBuffer) Reset() {
	b.buff.Reset()
	b.len = 0
}

func (b *stringBuffer) WriteInto(w *stringBuffer) {
	_, _ = b.buff.WriteTo(&w.buff)
	w.len += b.len
	b.len = 0
}

func (b *stringBuffer) String() string {
	return b.buff.String()
}

func (b *stringBuffer) WriteRune(r rune) {
	b.buff.WriteRune(r)
	b.len++
}

func (b *stringBuffer) WriteString(s string) {
	_, _ = b.buff.WriteString(s)
	b.len += len(s)
}
