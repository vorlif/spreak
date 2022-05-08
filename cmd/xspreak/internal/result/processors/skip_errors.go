package processors

import (
	"github.com/vorlif/spreak/xspreak/internal/result"
)

type skipErrors struct{}

func NewSkipErrors() Processor {
	return &skipErrors{}
}

func (s skipErrors) Process(issues []result.Issue) ([]result.Issue, error) {
	return filterIssues(issues, func(i *result.Issue) bool { return i.FromExtractor != "error_extractor" }), nil
}

func (s skipErrors) Name() string {
	return "skip_errors"
}
