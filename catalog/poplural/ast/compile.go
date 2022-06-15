package ast

import (
	"fmt"
	"strconv"
)

func CompileToString(forms *Forms) string {
	return fmt.Sprintf("nplurals=%d; plural=%s;", forms.NPlurals, compileNode(forms.Root))
}

func compileNode(node Node) string {
	switch v := node.(type) {
	case *ValueExpr:
		return strconv.FormatInt(v.Value, 10)
	case *OperandExpr:
		return "n"
	case *QuestionMarkExpr:
		return fmt.Sprintf("(%s) ? %s : %s", compileNode(v.Cond), compileNode(v.T), compileNode(v.F))
	case *BinaryExpr:
		switch v.Type() {
		case LogicalAnd:
			return compileNode(v.X) + " && " + compileNode(v.Y)
		case LogicalOr:
			return fmt.Sprintf("(%s || %s)", compileNode(v.X), compileNode(v.Y))
		case Equal:
			return compileNode(v.X) + " == " + compileNode(v.Y)
		case NotEqual:
			return compileNode(v.X) + " != " + compileNode(v.Y)
		case Greater:
			return compileNode(v.X) + " > " + compileNode(v.Y)
		case GreaterOrEqual:
			return compileNode(v.X) + " >= " + compileNode(v.Y)
		case Less:
			return compileNode(v.X) + " < " + compileNode(v.Y)
		case LessOrEqual:
			return compileNode(v.X) + " <= " + compileNode(v.Y)
		case Reminder:
			return compileNode(v.X) + " % " + compileNode(v.Y)
		}
	}
	return ""
}
