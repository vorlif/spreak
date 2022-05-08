package goextractors

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"strings"
)

func objToKey(obj types.Object) string {
	return fmt.Sprintf("%s.%s", obj.Pkg().Path(), obj.Name())
}

func searchSelector(expr interface{}) *ast.SelectorExpr {
	switch v := expr.(type) {
	case *ast.SelectorExpr:
		return v
	case *ast.Ident:
		if v.Obj == nil {
			break
		}
		return searchSelector(v.Obj.Decl)
	case *ast.ValueSpec:
		return searchSelector(v.Type)
	case *ast.Field:
		return searchSelector(v.Type)
	}
	return nil
}

func ExtractStringLiteral(expr ast.Expr) (string, ast.Node) {
	stack := []ast.Expr{expr}
	var b strings.Builder

	for len(stack) != 0 {
		n := len(stack) - 1
		elem := stack[n]
		stack = stack[:n]

		switch v := elem.(type) {
		//  Simple string with quotes or backqoutes
		case *ast.BasicLit:
			if v.Kind != token.STRING {
				return "", v
			}

			if unqouted, err := strconv.Unquote(v.Value); err != nil {
				b.WriteString(v.Value)
			} else {
				b.WriteString(unqouted)
			}
		// Concatenation of several string literals
		case *ast.BinaryExpr:
			if v.Op != token.ADD {
				return "", v
			}
			stack = append(stack, v.Y, v.X)
		case *ast.Ident:
			if v.Obj == nil {
				return "", v
			}
			switch z := v.Obj.Decl.(type) {
			case *ast.ValueSpec:
				if len(z.Values) == 0 {
					return "", v
				}
				stack = append(stack, z.Values[0])
			}
		default:
			return "", nil
		}
	}

	result := b.String()
	return result, expr
}

func calculatePosIdx(first, second int) int {
	if first > 0 && second > 0 {
		return first * second
	} else if first > 0 {
		return first
	} else {
		return second
	}
}
