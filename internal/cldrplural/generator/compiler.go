package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vorlif/spreak/internal/cldrplural/ast"
)

func compileNode(node ast.Node) string {
	switch e := node.(type) {
	case *ast.BinaryExpr:
		left := compileNode(e.X)
		right := compileNode(e.Y)

		if e.Op == ast.And {
			return left + " && " + right
		}
		return left + " || " + right
	case *ast.ModuloExpr:
		if e.Op.Operand == "n" {
			return fmt.Sprintf("math.Mod(%s, %d)", compileNode(e.Op), e.Value)
		}

		return fmt.Sprintf("%s %% %d", compileNode(e.Op), e.Value)
	case *ast.OperandExpr:
		return goOperandName(e.Operand)
	case *ast.InRelationExpr:
		var op *ast.OperandExpr
		if modNode, ok := e.X.(*ast.ModuloExpr); ok {
			op = modNode.Op
		} else {
			op = e.X.(*ast.OperandExpr)
		}

		var singleValues []string
		var valueRanges []string

		ast.Inspect(e.Y, func(node ast.Node) bool {
			switch valNode := node.(type) {
			case *ast.ValueExpr:
				singleValues = append(singleValues, strconv.FormatInt(valNode.Value, 10))
			case *ast.RangeExpr:
				valueRanges = append(valueRanges,
					strconv.FormatInt(valNode.From, 10),
					strconv.FormatInt(valNode.To, 10))
			}

			return true
		})

		compiledExpr := compileNode(e.X)
		isFloatOperand := op.Operand == "n"

		var b strings.Builder
		if len(valueRanges) > 0 {
			if isFloatOperand {
				b.WriteString(fmt.Sprintf("isFloatInRange(%s, ", compiledExpr))
			} else {
				b.WriteString(fmt.Sprintf("isIntInRange(%s, ", compiledExpr))
			}

			b.WriteString(strings.Join(valueRanges, ","))
			b.WriteString(")")
		}

		if len(singleValues) > 0 {
			if len(valueRanges) > 0 {
				b.WriteString(" || ")
			}

			if len(singleValues) == 1 {
				b.WriteString(compiledExpr)
				if len(valueRanges) == 0 && e.Op == ast.NotEqual {
					b.WriteString(" != ")
					b.WriteString(singleValues[0])
					return b.String()
				}

				b.WriteString(" == ")
				b.WriteString(singleValues[0])
			} else {
				if isFloatOperand {
					b.WriteString(fmt.Sprintf("isFloatOneOf(%s, ", compiledExpr))
				} else {
					b.WriteString(fmt.Sprintf("isIntOneOf(%s, ", compiledExpr))
				}

				b.WriteString(strings.Join(singleValues, ","))
				b.WriteString(")")
			}
		}

		if e.Op == ast.Equal {
			// If there are multiple checks they are separated by an "or" and should be bracketed.
			if len(singleValues)+(len(valueRanges)/2) > 1 {
				return fmt.Sprintf("( %s )", b.String())
			}

			return b.String()
		}

		return fmt.Sprintf("!( %s )", b.String())
	default:
		return ""
	}
}

func goOperandName(o string) string {
	switch o {
	case "n":
		return "ops.N"
	case "i":
		return "ops.I"
	case "v":
		return "ops.V"
	case "w":
		return "ops.W"
	case "f":
		return "ops.F"
	case "t":
		return "ops.T"
	case "c", "e":
		return "ops.C"
	default:
		panic("invalid operand " + o)
	}
}
