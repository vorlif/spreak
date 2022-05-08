package goextractors

import (
	"context"
	"go/ast"

	"github.com/vorlif/spreak/xspreak/internal/extract/extractors"
	"github.com/vorlif/spreak/xspreak/internal/result"
)

type globalAssignExtractor struct{}

func NewGlobalAssignExtractor() extractors.Extractor {
	return &globalAssignExtractor{}
}

func (v globalAssignExtractor) Run(_ context.Context, extractCtx *extractors.Context) ([]result.Issue, error) {
	var issues []result.Issue

	extractCtx.Inspector.WithStack([]ast.Node{&ast.ValueSpec{}}, func(rawNode ast.Node, push bool, stack []ast.Node) (proceed bool) {
		proceed = true
		if !push {
			return
		}

		node := rawNode.(*ast.ValueSpec)

		selector := searchSelector(node.Type)
		if selector == nil {
			return
		}

		tok := extractCtx.GetLocalizeTypeToken(selector)
		if tok != extractors.TypeSingular {
			return
		}

		in, ok := selector.X.(*ast.Ident)
		if !ok {
			return
		}

		pkg, _ := extractCtx.GetType(in)
		if pkg == nil {
			return
		}

		for _, value := range node.Values {
			msgID, stringNode := ExtractStringLiteral(value)
			if msgID == "" {
				return
			}

			issue := result.Issue{
				FromExtractor: v.Name(),
				MsgID:         msgID,
				Pkg:           pkg,
				CommentGroups: extractCtx.GetComments(pkg, stringNode, stack),
				Pos:           extractCtx.GetPosition(value.Pos()),
			}

			issues = append(issues, issue)
		}

		return
	})

	return issues, nil
}

func (v globalAssignExtractor) Name() string {
	return "consts_extractor"
}
