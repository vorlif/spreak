package mo

type File struct {
	Header   Header
	Messages []Message
}

type Header struct {
	ProjectIDVersion        string // Project-Id-Version: PACKAGE VERSION
	ReportMsgidBugsTo       string // Report-Msgid-Bugs-To: FIRST AUTHOR <EMAIL@ADDRESS>
	POTCreationDate         string // POT-Creation-Date: YEAR-MO-DA HO:MI+ZONE
	PORevisionDate          string // PO-Revision-Date: YEAR-MO-DA HO:MI+ZONE
	LastTranslator          string // Last-Translator: FIRST AUTHOR <EMAIL@ADDRESS>
	LanguageTeam            string // Language-Team:
	Language                string // Language: de
	MimeVersion             string // MIME-Version: 1.0
	ContentType             string // Content-Type: text/plain; charset=UTF-8
	ContentTransferEncoding string // Content-Transfer-Encoding: 8bit
	PluralForms             string // Plural-Forms: nplurals=2; plural=(n != 1);
	XGenerator              string // X-Generator: Poedit 3.0.1
	UnknownFields           map[string]string
}

type Message struct {
	Context   string
	ID        string
	IDPlural  string
	Str       string
	StrPlural []string
}
