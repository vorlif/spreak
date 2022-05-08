package processors

import (
	"github.com/vorlif/spreak/xspreak/internal/result"
)

type Processor interface {
	Process(issues []result.Issue) ([]result.Issue, error)
	Name() string
}
