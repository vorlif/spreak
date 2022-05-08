package processors

import (
	"github.com/vorlif/spreak/xspreak/internal/result"
)

type skipIgnoreFlag struct{}

func NewSkipIgnore() Processor {
	return &skipIgnoreFlag{}
}

func (s skipIgnoreFlag) Process(issues []result.Issue) ([]result.Issue, error) {
	return filterIssues(issues, func(i *result.Issue) bool {
		for _, f := range i.Flags {
			if f == "ignore" {
				return false
			}
		}

		return true
	}), nil
}

func (s skipIgnoreFlag) Name() string {
	return "skip_errors"
}
