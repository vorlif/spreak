package processors

import (
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/vorlif/spreak/internal/po"

	"github.com/vorlif/spreak/xspreak/internal/config"
	"github.com/vorlif/spreak/xspreak/internal/result"
)

type translationBuilder struct {
	cfg *config.Config
}

func BuildTranslations(cfg *config.Config) Processor {
	return &translationBuilder{
		cfg: cfg,
	}
}

func (s translationBuilder) Process(inIssues []result.Issue) ([]result.Issue, error) {
	outIssues := make([]result.Issue, 0, len(inIssues))

	absOut, errA := filepath.Abs(s.cfg.OutputDir)
	if errA != nil {
		absOut = s.cfg.OutputDir
	}

	for _, iss := range inIssues {
		var codeReferences []*po.Reference
		if !s.cfg.WriteNoLocation {
			path, errP := filepath.Rel(absOut, iss.Pos.Filename)
			if errP != nil {
				logrus.WithError(errP).Warn("Relative path could not be created")
				path = iss.Pos.Filename
			}

			ref := &po.Reference{
				Path:   path,
				Line:   iss.Pos.Line,
				Column: iss.Pos.Column,
			}
			codeReferences = append(codeReferences, ref)
		}

		iss.Message = &po.Message{
			Comment: &po.Comment{
				Extracted:  strings.Join(iss.Comment, "\n"),
				References: codeReferences,
				Flags:      iss.Flags,
			},
			Context:  iss.Context,
			ID:       iss.MsgID,
			IDPlural: iss.PluralID,
		}

		outIssues = append(outIssues, iss)
	}

	return outIssues, nil
}

func (s translationBuilder) Name() string {
	return "build_translation"
}
