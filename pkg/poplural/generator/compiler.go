package main

import (
	"strconv"
	"strings"

	"github.com/vorlif/spreak/pkg/poplural/ast"
)

func compileForms(forms *ast.Forms) string {
	if forms.NPlurals <= 1 {
		return "return 0"
	}

	var b strings.Builder

	if forms.NPlurals == 2 {
		b.WriteString("if ")
		compileBinaryPluralNode(forms.Root, &b)
		b.WriteString("{\nreturn 1\n}\n")
		b.WriteString("return 0")
	} else {
		compileNode(forms.Root, &b)
	}

	return b.String()
}

func compileBinaryPluralNode(node ast.Node, b *strings.Builder) {
	switch v := node.(type) {
	case *ast.ValueExpr:
		b.WriteString(strconv.FormatInt(v.Value, 10))
	case *ast.OperandExpr:
		b.WriteString("n")
	case *ast.BinaryExpr:
		switch v.Type() {
		case ast.LogicalAnd:
			compileBinaryPluralNode(v.X, b)
			b.WriteString(" && ")
			compileBinaryPluralNode(v.Y, b)
		case ast.LogicalOr:
			b.WriteString("(")
			compileBinaryPluralNode(v.X, b)
			b.WriteString(" || ")
			compileBinaryPluralNode(v.Y, b)
			b.WriteString(")")
		case ast.Equal:
			compileBinaryPluralNode(v.X, b)
			b.WriteString(" == ")
			compileBinaryPluralNode(v.Y, b)
		case ast.NotEqual:
			compileBinaryPluralNode(v.X, b)
			b.WriteString(" != ")
			compileBinaryPluralNode(v.Y, b)
		case ast.Greater:
			compileBinaryPluralNode(v.X, b)
			b.WriteString(" > ")
			compileBinaryPluralNode(v.Y, b)
		case ast.GreaterOrEqual:
			compileBinaryPluralNode(v.X, b)
			b.WriteString(" >= ")
			compileBinaryPluralNode(v.Y, b)
		case ast.Less:
			compileBinaryPluralNode(v.X, b)
			b.WriteString(" < ")
			compileBinaryPluralNode(v.Y, b)
		case ast.LessOrEqual:
			compileBinaryPluralNode(v.X, b)
			b.WriteString(" <= ")
			compileBinaryPluralNode(v.Y, b)
		case ast.Reminder:
			compileBinaryPluralNode(v.X, b)
			b.WriteString("%")
			compileBinaryPluralNode(v.Y, b)
		}
	}
}

// compileNode converts a rule tree to go code.
func compileNode(node ast.Node, b *strings.Builder) {
	switch v := node.(type) {
	case *ast.QuestionMarkExpr:
		b.WriteString("if ")
		compileBinaryPluralNode(v.Cond, b)
		if vn, ok := v.T.(*ast.ValueExpr); ok {
			b.WriteString("{\n")
			b.WriteString("return ")
			b.WriteString(strconv.FormatInt(vn.Value, 10))
			b.WriteString("\n}\n")
		} else {
			compileNode(v.T, b)
		}

		if vn, ok := v.F.(*ast.ValueExpr); ok {
			b.WriteString("return ")
			b.WriteString(strconv.FormatInt(vn.Value, 10))
		} else {
			compileNode(v.F, b)
		}
	default:
		compileBinaryPluralNode(v, b)
	}
}
