package cldrplural

import (
	"math"

	"github.com/vorlif/spreak/catalog/cldrplural/ast"
	"github.com/vorlif/spreak/internal/util"
)

func MustParseRules(rules map[Category]string) *RuleSet {
	f, err := ParseRules(rules)
	if err != nil {
		panic(err)
	}

	return f
}

// ParseRules creates a RuleSet from a set of rules, which can be evaluated at runtime.
func ParseRules(rawRules map[Category]string) (*RuleSet, error) {
	categories := make([]Category, 0, len(rawRules))
	rules := make(map[Category]*ast.Rule)

	for cat, rawRule := range rawRules {
		rule, err := ast.Parse(rawRule)
		if err != nil {
			return nil, err
		}
		rules[cat] = rule
		categories = append(categories, cat)
	}

	if _, hasOther := rawRules[Other]; !hasOther {
		categories = append(categories, Other)
	}

	formF := func(rules map[Category]*ast.Rule) FormFunc {
		return func(ops *Operands) Category {
			for cat, rule := range rules {
				if evaluate(rule, ops) {
					return cat
				}
			}

			return Other
		}
	}(rules)

	return &RuleSet{Categories: categories, FormFunc: formF}, nil
}

// Evaluates whether an abstract rule applies to an operation.
func evaluate(rule *ast.Rule, ops *Operands) bool {
	// Other has no conditions and therefore no tree
	if rule.Root == nil {
		return true
	}
	return evaluateNode(rule.Root, ops)
}

func evaluateNode(node ast.Node, ops *Operands) bool {
	switch e := node.(type) {
	case *ast.BinaryExpr:
		left := evaluateNode(e.X, ops)
		right := evaluateNode(e.Y, ops)
		if e.Op == ast.And {
			return left && right
		}
		return left || right
	case *ast.InRelationExpr:
		evaluatedValue := evaluateExpression(e.X, ops)
		inRelation := false

		ast.Inspect(e.Y, func(node ast.Node) bool {
			switch valNode := node.(type) {
			case *ast.ValueExpr:
				if float64(valNode.Value) == evaluatedValue {
					inRelation = true
					return false
				}
			case *ast.RangeExpr:
				if isFloatInRange(evaluatedValue, valNode.From, valNode.To) {
					inRelation = true
					return false
				}
			}

			return true
		})

		if e.Op == ast.NotEqual {
			return !inRelation
		}

		return inRelation
	}

	return false
}

func evaluateExpression(node ast.Node, ops *Operands) float64 {
	switch e := node.(type) {
	case *ast.ModuloExpr:
		val := e.Value
		switch e.Op.Operand {
		case "n":
			return math.Mod(ops.N, float64(val))
		case "i":
			return float64(ops.I % val)
		case "v":
			return float64(ops.V % val)
		case "w":
			return float64(ops.W % val)
		case "f":
			return float64(ops.F % val)
		case "t":
			return float64(ops.T % val)
		case "c", "e":
			return float64(ops.C % val)
		default:
			panic("invalid operand " + e.Op.Operand)
		}
	case *ast.OperandExpr:
		switch e.Operand {
		case "n":
			return ops.N
		case "i":
			return float64(ops.I)
		case "v":
			return float64(ops.V)
		case "w":
			return float64(ops.W)
		case "f":
			return float64(ops.F)
		case "t":
			return float64(ops.T)
		case "c", "e":
			return float64(ops.C)
		default:
			panic("invalid operand " + e.Operand)
		}
	default:
		panic("not an expression")
	}
}

// isFloatInRange tests whether a float64 value x is within a range [from, to].
// Multiple ranges can be specified.
func isFloatInRange(x float64, rangeValues ...int64) bool {
	for i := 0; i < len(rangeValues); i += 2 {
		lower := rangeValues[i]
		upper := rangeValues[i+1]

		// If the value is not within the range, not every value needs to be checked.
		if !(x >= float64(lower) && x <= float64(upper)) {
			continue
		}

		// Checking each value
		for v := lower; v <= upper; v++ {
			if util.FloatEqual(float64(v), x) {
				return true
			}
		}
	}

	return false
}

func isFloatOneOf(target float64, vals ...int64) bool {
	for _, val := range vals {
		if util.FloatEqual(float64(val), target) {
			return true
		}
	}

	return false
}
