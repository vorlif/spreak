package goextractors

import (
	"context"
	"go/ast"

	"github.com/vorlif/spreak/xspreak/internal/extract/extractors"
	"github.com/vorlif/spreak/xspreak/internal/result"
)

type funcCallExtractor struct{}

func NewFuncCallExtractor() extractors.Extractor {
	return &funcCallExtractor{}
}

func (v funcCallExtractor) Run(_ context.Context, extractCtx *extractors.Context) ([]result.Issue, error) {
	var issues []result.Issue

	extractCtx.Inspector.WithStack([]ast.Node{&ast.CallExpr{}}, func(rawNode ast.Node, push bool, stack []ast.Node) (proceed bool) {
		proceed = true
		if !push {
			return
		}

		node := rawNode.(*ast.CallExpr)
		if len(node.Args) == 0 {
			return
		}

		var ident *ast.Ident
		if selector := searchSelector(node.Fun); selector != nil {
			ident = selector.Sel
		} else {
			idt, ok := node.Fun.(*ast.Ident)
			if !ok {
				return
			}
			ident = idt
		}

		pkg, obj := extractCtx.GetType(ident)
		if pkg == nil {
			return
		}

		if tok := extractCtx.GetLocalizeTypeToken(ident); tok == extractors.TypeSingular {
			raw, stringNode := ExtractStringLiteral(node.Args[0])
			if raw == "" {
				return
			}

			issue := result.Issue{
				FromExtractor: v.Name(),
				MsgID:         raw,
				Pkg:           pkg,
				CommentGroups: extractCtx.GetComments(pkg, stringNode, stack),
				Pos:           extractCtx.GetPosition(node.Args[0].Pos()),
			}

			issues = append(issues, issue)
		}

		funcParameterDefs := extractCtx.Definitions.GetFields(objToKey(obj))
		if funcParameterDefs == nil {
			return
		}

		issue := result.Issue{
			FromExtractor: v.Name(),
			Pkg:           pkg,
			Pos:           extractCtx.GetPosition(node.Args[0].Pos()),
			CommentGroups: extractCtx.GetComments(pkg, node.Args[0], stack),
		}
		for _, def := range funcParameterDefs {
			for i, arg := range node.Args {
				if def.FieldPos != i {
					continue
				}

				raw, _ := ExtractStringLiteral(arg)
				if raw == "" {
					return
				}
				switch def.Token {
				case extractors.TypeSingular:
					issue.MsgID = raw
				case extractors.TypePlural:
					issue.PluralID = raw
				case extractors.TypeContext:
					issue.Context = raw
				case extractors.TypeDomain:
					issue.Domain = raw
				}
			}
		}

		issues = append(issues, issue)
		return
	})

	return issues, nil
}

func (v funcCallExtractor) Name() string {
	return "funccall_extractor"
}
