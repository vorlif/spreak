package goextractors

import (
	"context"
	"go/ast"

	"github.com/vorlif/spreak/xspreak/internal/result"

	"github.com/vorlif/spreak/xspreak/internal/config"
	"github.com/vorlif/spreak/xspreak/internal/extract/extractors"
)

type errorExtractor struct{}

func NewErrorExtractor() extractors.Extractor {
	return &errorExtractor{}
}

func (v errorExtractor) Run(_ context.Context, extractCtx *extractors.Context) ([]result.Issue, error) {
	var issues []result.Issue

	extractCtx.Inspector.WithStack([]ast.Node{&ast.CallExpr{}}, func(rawNode ast.Node, push bool, stack []ast.Node) (proceed bool) {
		proceed = true
		if !push {
			return
		}

		node := rawNode.(*ast.CallExpr)
		if len(node.Args) != 1 {
			return
		}

		selector := searchSelector(node.Fun)
		if selector == nil {
			return
		}

		pkg, obj := extractCtx.GetType(selector.Sel)
		if pkg == nil {
			return
		}

		if "errors" != obj.Pkg().Path() || !config.ShouldExtractPackage(pkg.PkgPath) {
			return
		}

		msgID, _ := ExtractStringLiteral(node.Args[0])
		if msgID == "" {
			return
		}

		issue := result.Issue{
			FromExtractor: v.Name(),
			MsgID:         msgID,
			Pkg:           pkg,
			Context:       extractCtx.Config.ErrorContext,
			CommentGroups: extractCtx.GetComments(pkg, node, stack),
			Pos:           extractCtx.GetPosition(node.Args[0].Pos()),
		}

		issues = append(issues, issue)

		return
	})

	return issues, nil
}

func (v errorExtractor) Name() string {
	return "error_extractor"
}
