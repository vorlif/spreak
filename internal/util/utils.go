package util

import (
	"bytes"
	"strings"
	"unicode"
)

func DecodePoString(text string) string {
	lines := strings.Split(text, "\n")
	for idx := 0; idx < len(lines); idx++ {
		left := strings.Index(lines[idx], `"`)
		right := strings.LastIndex(lines[idx], `"`)
		if left < 0 || right < 0 || left == right {
			lines[idx] = ""
			continue
		}
		line := lines[idx][left+1 : right]

		data := make([]byte, 0, len(line))
		for i := 0; i < len(line); i++ {
			if line[i] != '\\' {
				data = append(data, line[i])
				continue
			}
			if i+1 >= len(line) {
				break
			}
			switch line[i+1] {
			case 'n': // \\n -> \n
				data = append(data, '\n')
				i++
			case 't': // \\t -> \n
				data = append(data, '\t')
				i++
			case '\\': // \\\ -> ?
				data = append(data, '\\')
				i++
			}
		}
		lines[idx] = string(data)
	}
	return strings.Join(lines, "")
}

func EncodePoString(s string, lim int) string {

	var lines []string
	if lim <= 0 {
		lines = encodePoStringWithoutWrap(s)
	} else {
		lines = encodePoStringWithWrap(s, lim)
	}

	// Single line msgid / entry starts with msgid "text"
	if len(lines) == 1 {
		// Single line
		return lines[0]
	} else if len(lines) == 2 && (lines[1] == "" || lines[1] == `""`) {
		// Single line with newline at end
		return lines[0]
	}

	lines = append(lines, "")
	copy(lines[1:], lines)
	lines[0] = `""`

	lastIdx := len(lines) - 1
	if lines[lastIdx] == "" || lines[lastIdx] == `""` {
		lines = lines[:lastIdx]
	}

	return strings.Join(lines, "\n")
}

func encodePoStringWithWrap(s string, pageWidth int) []string {
	lines := make([]string, 0, 2) // The average message is two lines long

	currentLine := &stringBuffer{}
	var lastWordBuf stringBuffer
	var currentSpaceBuf stringBuffer

	currentLine.WriteRune('"')
	var currentLen int

	createNewLine := func() {
		lines = append(lines, currentLine.String())
		currentLine.Reset()
		currentLine.WriteRune('"')
		currentLen = 0
	}

	for _, char := range s {
		currentLen = currentLine.Len() - 1
		if char == '\n' {
			lastWordBuf.WriteInto(currentLine)
			currentSpaceBuf.WriteInto(currentLine)

			currentLine.WriteString(`\n"`)
			createNewLine()

		} else if unicode.IsSpace(char) && char != nbsp {
			if currentLen+lastWordBuf.Len()+currentSpaceBuf.Len()-1 >= pageWidth && currentLen > 1 {
				currentLine.WriteRune('"')
				createNewLine()
			}

			currentSpaceBuf.WriteRune(char)

			if currentLen+lastWordBuf.Len()+currentSpaceBuf.Len()-1 >= pageWidth && currentLen > 1 {
				currentLine.WriteRune('"')
				createNewLine()
			}
		} else {
			if currentSpaceBuf.Len() > 0 {
				if currentLen+lastWordBuf.Len()+currentSpaceBuf.Len()+1 >= pageWidth && currentLen == 0 {
					lastWordBuf.WriteInto(currentLine)
					currentSpaceBuf.WriteInto(currentLine)
					currentLine.WriteRune('"')
					createNewLine()
				}

				lastWordBuf.WriteInto(currentLine)
				currentSpaceBuf.WriteInto(currentLine)
			}

			lastWordBuf.WriteRune(char)
		}
	}

	lastWordBuf.WriteInto(currentLine)
	currentSpaceBuf.WriteInto(currentLine)
	currentLine.WriteRune('"')

	if remain := currentLine.String(); remain != `""` || len(lines) == 0 {
		lines = append(lines, remain)
	}

	return lines
}

func encodePoStringWithoutWrap(s string) []string {
	lines := make([]string, 0, 2) // The average message is two lines long

	init := make([]byte, 0, len(s))
	buf := bytes.NewBuffer(init)
	buf.WriteRune('"')
	for _, char := range s {
		// A newline closes the line
		if char == '\n' {
			buf.WriteString(`\n"`)
			lines = append(lines, buf.String())
			buf.Reset()
			buf.WriteRune('"')
			continue
		}

		buf.WriteRune(char)
	}
	buf.WriteRune('"')

	if remain := buf.String(); remain != `""` || len(lines) == 0 {
		lines = append(lines, remain)
	}

	return lines
}
