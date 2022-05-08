package goextractors

import (
	"context"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"

	"github.com/vorlif/spreak/xspreak/internal/extract/extractors"
	"github.com/vorlif/spreak/xspreak/internal/result"
)

type structDefExtractor struct{}

func NewStructDefExtractor() extractors.Extractor {
	return &structDefExtractor{}
}

func (v structDefExtractor) Run(_ context.Context, extractCtx *extractors.Context) ([]result.Issue, error) {
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

		var obj types.Object
		var pkg *packages.Package
		switch val := node.Type.(type) {
		case *ast.SelectorExpr:
			pkg, obj = extractCtx.GetType(val.Sel)
			if pkg == nil {
				return
			}
		case *ast.Ident:
			pkg, obj = extractCtx.GetType(val)
			if pkg == nil {
				return
			}
		default:
			return
		}

		if structAttr := extractCtx.Definitions.GetFields(objToKey(obj)); structAttr == nil {
			return
		}

		structIssues := extractStruct(extractCtx, node, obj, pkg, stack)
		issues = append(issues, structIssues...)

		return
	})

	return issues, nil
}

func (v structDefExtractor) Name() string {
	return "struct_extractor"
}

func extractStruct(extractCtx *extractors.Context, node *ast.CompositeLit, obj types.Object, pkg *packages.Package, stack []ast.Node) []result.Issue {
	var issues []result.Issue
	issue := result.Issue{
		Pkg:           pkg,
		Pos:           extractCtx.GetPosition(node.Pos()),
		CommentGroups: extractCtx.GetComments(pkg, node, stack),
	}
	definitionKey := objToKey(obj)
	if _, isKv := node.Elts[0].(*ast.KeyValueExpr); isKv {
		for _, elt := range node.Elts {
			kve, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}

			idt, ok := kve.Key.(*ast.Ident)
			if !ok {
				continue
			}

			def := extractCtx.Definitions.Get(definitionKey, idt.Name)
			if def == nil {
				continue
			}

			raw, stringNode := ExtractStringLiteral(kve.Value)
			if raw == "" {
				continue
			}

			if def.Token == extractors.TypeSingular && issue.MsgID != "" {
				issues = append(issues, issue)
				issue = result.Issue{
					Pkg:           pkg,
					Pos:           extractCtx.GetPosition(node.Pos()),
					CommentGroups: extractCtx.GetComments(pkg, stringNode, stack),
				}
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
	} else {
		for _, attrDef := range extractCtx.Definitions.GetFields(definitionKey) {
			for i, elt := range node.Elts {
				if attrDef.FieldPos != i {
					continue
				}

				raw, stringNode := ExtractStringLiteral(elt)
				if raw == "" {
					continue
				}

				if attrDef.Token == extractors.TypeSingular && issue.MsgID != "" {
					issues = append(issues, issue)
					issue = result.Issue{
						Pkg:           pkg,
						Pos:           extractCtx.GetPosition(node.Pos()),
						CommentGroups: extractCtx.GetComments(pkg, stringNode, stack),
					}
				}

				switch attrDef.Token {
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
	}

	issues = append(issues, issue)
	return issues
}
