package po

import (
	"fmt"
	"io"

	"github.com/vorlif/spreak/internal/util"
)

type Header struct {
	Comment                 *util.Comment // Header Comments
	ProjectIDVersion        string        // Project-Id-Version: PACKAGE VERSION
	ReportMsgidBugsTo       string        // Report-Msgid-Bugs-To: FIRST AUTHOR <EMAIL@ADDRESS>
	POTCreationDate         string        // POT-Creation-Date: YEAR-MO-DA HO:MI+ZONE
	PORevisionDate          string        // PO-Revision-Date: YEAR-MO-DA HO:MI+ZONE
	LastTranslator          string        // Last-Translator: FIRST AUTHOR <EMAIL@ADDRESS>
	LanguageTeam            string        // Language-Team:
	Language                string        // Language: de
	MimeVersion             string        // MIME-Version: 1.0
	ContentType             string        // Content-Type: text/plain; charset=UTF-8
	ContentTransferEncoding string        // Content-Transfer-Encoding: 8bit
	PluralForms             string        // Plural-Forms: nplurals=2; plural=(n != 1);
	XGenerator              string        // X-Generator: Poedit 3.0.1
	UnknownFields           map[string]string
}

func (p *Header) WriteTo(w io.StringWriter, wrapWidth int) error {
	if p.Comment != nil {
		if err := p.Comment.WriteTo(w, wrapWidth); err != nil {
			return err
		}
	}
	lines := []string{
		`msgid ""` + "\n",
		`msgstr ""` + "\n",
		fmt.Sprintf(`"%s: %s\n"`+"\n", util.HeaderProjectIDVersion, p.ProjectIDVersion),
		fmt.Sprintf(`"%s: %s\n"`+"\n", util.HeaderReportMsgIDBugsTo, p.ReportMsgidBugsTo),
		fmt.Sprintf(`"%s: %s\n"`+"\n", util.HeaderPOTCreationDate, p.POTCreationDate),
		fmt.Sprintf(`"%s: %s\n"`+"\n", util.HeaderPORevisionDate, p.PORevisionDate),
		fmt.Sprintf(`"%s: %s\n"`+"\n", util.HeaderLastTranslator, p.LastTranslator),
		fmt.Sprintf(`"%s: %s\n"`+"\n", util.HeaderLanguageTeam, p.LanguageTeam),
		fmt.Sprintf(`"%s: %s\n"`+"\n", util.HeaderLanguage, p.Language),
	}

	if p.MimeVersion != "" {
		lines = append(lines, fmt.Sprintf(`"%s: %s\n"`+"\n", util.HeaderMIMEVersion, p.MimeVersion))
	}

	lines = append(lines,
		fmt.Sprintf(`"%s: %s\n"`+"\n", util.HeaderContentType, p.ContentType),
		fmt.Sprintf(`"%s: %s\n"`+"\n", util.HeaderContentTransferEncoding, p.ContentTransferEncoding),
	)

	if p.PluralForms != "" {
		lines = append(lines, fmt.Sprintf(`"%s: %s\n"`+"\n", util.HeaderPluralForms, p.PluralForms))
	}

	if p.XGenerator != "" {
		lines = append(lines, fmt.Sprintf(`"%s: %s\n"`+"\n", util.HeaderXGenerator, p.XGenerator))
	}
	for k, v := range p.UnknownFields {
		lines = append(lines, fmt.Sprintf(`"%s: %s\n"`+"\n", k, v))
	}

	for _, line := range lines {
		if _, err := w.WriteString(line); err != nil {
			return err
		}
	}
	return nil
}
