package extract

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/packages"

	"github.com/vorlif/spreak/xspreak/internal/config"
	"github.com/vorlif/spreak/xspreak/internal/extract/extractors"

	"github.com/vorlif/spreak/xspreak/internal/result"
	"github.com/vorlif/spreak/xspreak/internal/result/processors"
	"github.com/vorlif/spreak/xspreak/internal/util"
)

type Runner struct {
	Processors []processors.Processor
	Log        *logrus.Entry
}

func NewRunner(cfg *config.Config, pkgs map[string]*packages.Package) (*Runner, error) {
	p := []processors.Processor{
		processors.NewSkipEmptyMsgID(),
	}

	if !cfg.ExtractErrors {
		p = append(p, processors.NewSkipErrors())
	}

	p = append(p,
		processors.NewCommentCleaner(cfg),
		processors.NewSkipIgnore(),
		processors.BuildTranslations(cfg))

	ret := &Runner{
		Processors: p,
		Log:        logrus.WithField("service", "Runner"),
	}

	return ret, nil
}

func (r Runner) Run(ctx context.Context, extractCtx *extractors.Context, extractors []extractors.Extractor) ([]result.Issue, error) {
	defer util.TrackTime(time.Now(), "Extracting the issues")
	var issues []result.Issue
	for _, extr := range extractors {
		extractedIssues, err := extr.Run(ctx, extractCtx)
		if err != nil {
			r.Log.Warnf("Can't run extractor %s: %v", extr.Name(), err)
		} else {
			issues = append(issues, extractedIssues...)
		}
	}

	return r.processIssues(issues), nil
}

func (r *Runner) processIssues(issues []result.Issue) []result.Issue {
	defer util.TrackTime(time.Now(), "Process the issues")
	for _, p := range r.Processors {
		var newIssues []result.Issue
		var err error

		newIssues, err = p.Process(issues)
		if err != nil {
			r.Log.Warnf("Can't process result by %s processor: %s", p.Name(), err)
		} else {
			issues = newIssues
		}

		if issues == nil {
			issues = []result.Issue{}
		}
	}

	return issues
}
