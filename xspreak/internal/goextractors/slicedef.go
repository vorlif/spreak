package goextractors

import (
	"context"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"

	"github.com/vorlif/spreak/xspreak/internal/result"

	"github.com/vorlif/spreak/xspreak/internal/extract/extractors"
)

type sliceDefExtractor struct{}

func NewSliceDefExtractor() extractors.Extractor {
	return &sliceDefExtractor{}
}

func (v sliceDefExtractor) Run(_ context.Context, extractCtx *extractors.Context) ([]result.Issue, error) {
	var issues []result.Issue

	extractCtx.Inspector.WithStack([]ast.Node{&ast.CompositeLit{}}, func(rawNode ast.Node, push bool, stack []ast.Node) (proceed bool) {
		proceed = true
		if !push {
			return
		}

		node := rawNode.(*ast.CompositeLit)
		if len(node.Elts) == 0 {
			return
		}

		arrayTye, ok := node.Type.(*ast.ArrayType)
		if !ok {
			return
		}

		var obj types.Object
		var pkg *packages.Package
		var token extractors.TypeToken
		switch val := arrayTye.Elt.(type) {
		case *ast.SelectorExpr:
			pkg, obj = extractCtx.GetType(val.Sel)
			if pkg == nil {
				return
			}
			token = extractCtx.GetLocalizeTypeToken(val.Sel)
		case *ast.Ident:
			pkg, obj = extractCtx.GetType(val)
			if pkg == nil {
				return
			}
			token = extractCtx.GetLocalizeTypeToken(val)
		case *ast.StarExpr:
			switch pointerExpr := val.X.(type) {
			case *ast.SelectorExpr:
				pkg, obj = extractCtx.GetType(pointerExpr.Sel)
				if pkg == nil {
					return
				}
				token = extractCtx.GetLocalizeTypeToken(pointerExpr.Sel)
			case *ast.Ident:
				pkg, obj = extractCtx.GetType(pointerExpr)
				if pkg == nil {
					return
				}
				token = extractCtx.GetLocalizeTypeToken(pointerExpr)

			default:
				return
			}
		default:
			return
		}

		// Array of strings
		if token == extractors.TypeSingular {
			for _, elt := range node.Elts {
				msgID, stringNode := ExtractStringLiteral(elt)
				if msgID == "" {
					continue
				}

				issue := result.Issue{
					FromExtractor: v.Name(),
					MsgID:         msgID,
					Pkg:           pkg,
					CommentGroups: extractCtx.GetComments(pkg, stringNode, stack),
					Pos:           extractCtx.GetPosition(node.Pos()),
				}

				issues = append(issues, issue)
			}

			return
		}

		structAttr := extractCtx.Definitions.GetFields(objToKey(obj))
		if structAttr == nil {
			return
		}

		for _, elt := range node.Elts {
			compLit, isCompLit := elt.(*ast.CompositeLit)
			if !isCompLit {
				continue
			}

			structIssues := extractStruct(extractCtx, compLit, obj, pkg, stack)
			issues = append(issues, structIssues...)
		}

		return
	})

	return issues, nil
}

func (v sliceDefExtractor) Name() string {
	return "slicedef_extractor"
}
