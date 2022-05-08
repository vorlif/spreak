package encoder

import (
	"fmt"
	"io"
	"time"

	"github.com/vorlif/spreak/xspreak/internal/result"

	"github.com/vorlif/spreak/internal/po"
	"github.com/vorlif/spreak/internal/util"
	"github.com/vorlif/spreak/xspreak/internal/config"
)

type Encoder interface {
	Encode(issues []result.Issue) error
}

type potEncoder struct {
	cfg *config.Config
	w   io.StringWriter
}

func NewPotEncoder(cfg *config.Config, w io.StringWriter) Encoder {
	return &potEncoder{cfg: cfg, w: w}
}

func (e *potEncoder) Encode(issues []result.Issue) error {
	var header *po.Header

	if !e.cfg.OmitHeader {
		header = e.buildHeader()
	}

	file := &po.File{
		Header:   header,
		Messages: make(map[string]map[string]*util.Message),
	}

	for _, iss := range issues {
		file.AddMessage(iss.Translation)
	}

	return file.WriteTo(e.w, e.cfg.WrapWidth)
}

func (e *potEncoder) buildHeader() *po.Header {
	headerComment := fmt.Sprintf(`SOME DESCRIPTIVE TITLE.
Copyright (C) YEAR %s
This file is distributed under the same license as the %s package.
FIRST AUTHOR <EMAIL@ADDRESS>, YEAR.
`, e.cfg.CopyrightHolder, e.cfg.PackageName)
	return &po.Header{
		Comment: &util.Comment{
			Translator:     headerComment,
			Extracted:      "",
			References:     nil,
			Flags:          []string{"fuzzy"},
			PrevMsgContext: "",
			PrevMsgID:      "",
		},
		ProjectIDVersion:        e.cfg.PackageName,
		ReportMsgidBugsTo:       e.cfg.BugsAddress,
		POTCreationDate:         time.Now().Format("2006-01-02 15:04-0700"),
		PORevisionDate:          "YEAR-MO-DA HO:MI+ZONE",
		LastTranslator:          "FULL NAME <EMAIL@ADDRESS>",
		LanguageTeam:            "LANGUAGE <LL@li.org>",
		Language:                "",
		MimeVersion:             "1.0",
		ContentType:             "text/plain; charset=UTF-8",
		ContentTransferEncoding: "8bit",
		PluralForms:             "", // alternative  "nplurals=INTEGER; plural=EXPRESSION;"
	}
}
