package extractors

import (
	"context"

	"github.com/vorlif/spreak/xspreak/internal/result"
)

type Extractor interface {
	Run(ctx context.Context, extractCtx *Context) ([]result.Issue, error)
	Name() string
}
