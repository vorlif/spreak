package goextractors

import (
	"context"
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vorlif/spreak/xspreak/internal/config"
	"github.com/vorlif/spreak/xspreak/internal/extract"
	"github.com/vorlif/spreak/xspreak/internal/extract/extractors"
	"github.com/vorlif/spreak/xspreak/internal/result"
)

const (
	testdataDir = "../../../../testdata/xspreak/project"
)

func TestPrintAst(t *testing.T) {
	fset := token.NewFileSet() // positions are relative to fset
	_, err := parser.ParseFile(fset, filepath.Join(testdataDir, "slices.go"), nil, 0)
	if err != nil {
		panic(err)
	}

	// ast.Print(fset, f)
}

func runExtraction(t *testing.T, dir string, testExtractors ...extractors.Extractor) []result.Issue {
	cfg := config.NewDefault()
	cfg.SourceDir = dir
	cfg.ExtractErrors = true
	ctx := context.Background()
	contextLoader := extract.NewContextLoader(cfg)

	extractCtx, err := contextLoader.Load(ctx)
	require.NoError(t, err)

	runner, err := extract.NewRunner(cfg, extractCtx.Packages)
	require.NoError(t, err)

	e := []extractors.Extractor{NewDefinitionExtractor()}
	if len(testExtractors) > 0 {
		e = append(e, testExtractors...)
	}
	issues, err := runner.Run(ctx, extractCtx, e)
	require.NoError(t, err)
	return issues
}

func collectIssueStrings(issues []result.Issue) []string {
	collection := make([]string, 0, len(issues))
	for _, issue := range issues {
		collection = append(collection, issue.MsgID)
		if issue.PluralID != "" {
			collection = append(collection, issue.PluralID)
		}

		if issue.Context != "" {
			collection = append(collection, issue.Context)
		}

		if issue.Domain != "" {
			collection = append(collection, issue.Domain)
		}
	}
	return collection
}
