package goextractors

import (
	"context"
	"go/ast"

	"github.com/vorlif/spreak/xspreak/internal/extract/extractors"
	"github.com/vorlif/spreak/xspreak/internal/result"
)

type varAssignExtractor struct{}

func NewVariablesExtractor() extractors.Extractor {
	return &varAssignExtractor{}
}

func (v varAssignExtractor) Run(_ context.Context, extractCtx *extractors.Context) ([]result.Issue, error) {
	var issues []result.Issue

	extractCtx.Inspector.WithStack([]ast.Node{&ast.AssignStmt{}}, func(rawNode ast.Node, push bool, stack []ast.Node) (proceed bool) {
		proceed = true
		if !push {
			return
		}
		node := rawNode.(*ast.AssignStmt)
		if len(node.Lhs) == 0 || len(node.Rhs) == 0 {
			return
		}

		token, ident := extractCtx.SearchIdentAndToken(node.Lhs[0])
		if token == extractors.TypeNone {
			return
		}

		pkg, _ := extractCtx.GetType(ident)
		if pkg == nil {
			return
		}

		if token == extractors.TypeSingular {
			msgID, stringNode := ExtractStringLiteral(node.Rhs[0])
			if msgID == "" {
				return
			}

			issue := result.Issue{
				FromExtractor: v.Name(),
				MsgID:         msgID,
				Pkg:           pkg,
				CommentGroups: extractCtx.GetComments(pkg, stringNode, stack),
				Pos:           extractCtx.GetPosition(node.Rhs[0].Pos()),
			}

			issues = append(issues, issue)
			return
		}

		return
	})

	return issues, nil
}

func (v varAssignExtractor) Name() string {
	return "varassign_extractor"
}
