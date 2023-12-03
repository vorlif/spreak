package poplural

import (
	"math"

	"github.com/vorlif/spreak/catalog/poplural/ast"
	"github.com/vorlif/spreak/internal/util"
)

// Rule represents a rule as specified in the header of the po file under `Plural-Forms`.
//
// Example: nplurals=2; plural=n != 1;.
type Rule struct {
	// Specifies how many plural forms this rule supports.
	NPlurals int
	// Function that returns the appropriate plural form for a number.
	FormFunc func(n int64) int
}

// Evaluate returns the appropriate plural form for any value.
// If the value is not an integer, it is formatted accordingly.
// If formatting is not possible, an error is returned.
func (f *Rule) Evaluate(a any) (int, error) {
	num, err := util.ToNumber(a)
	if err != nil {
		return 0, err
	}

	num = math.RoundToEven(math.Abs(num))
	return f.FormFunc(int64(num)), nil
}

// MustParse is like Parse, but panics if the given rule string cannot be parsed.
// It simplifies safe initialization of Rule values.
func MustParse(raw string) *Rule {
	rule, err := Parse(raw)
	if err != nil {
		panic(err)
	}
	return rule
}

// Parse parses a plural forms header and returns a function to evaluate this header.
// If for a header there is already a predefined function, this function will be returned.
// If parsing failed it returns an error.
func Parse(raw string) (*Rule, error) {
	parsed, err := ast.Parse(raw)
	if err != nil {
		return nil, err
	}

	// Use of built-in functions, if available
	compiledRaw := ast.CompileToString(parsed)
	if rule := forRawRule(compiledRaw); rule != nil {
		return rule, nil
	}

	f := &Rule{
		NPlurals: parsed.NPlurals,
		FormFunc: generateFormFunc(parsed),
	}
	return f, nil
}

func generateFormFunc(forms *ast.Forms) func(n int64) int {
	if forms.Root == nil {
		return func(n int64) int { return 0 }
	}

	return func(n int64) int {
		return int(evaluateNode(forms.Root, n))
	}
}

// evaluateNode computes the plural form for a number and a rule that was parsed at runtime.
func evaluateNode(node ast.Node, num int64) int64 {
	var conditionTrue bool

	switch v := node.(type) {
	case *ast.ValueExpr:
		return v.Value
	case *ast.OperandExpr:
		return num
	case *ast.QuestionMarkExpr:
		if evaluateNode(v.Cond, num) == 1 {
			return evaluateNode(v.T, num)
		}
		return evaluateNode(v.F, num)
	case *ast.BinaryExpr:
		switch v.Type() {
		case ast.LogicalAnd:
			conditionTrue = evaluateNode(v.X, num) == 1 && evaluateNode(v.Y, num) == 1
		case ast.LogicalOr:
			conditionTrue = evaluateNode(v.X, num) == 1 || evaluateNode(v.Y, num) == 1
		case ast.Equal:
			conditionTrue = evaluateNode(v.X, num) == evaluateNode(v.Y, num)
		case ast.NotEqual:
			conditionTrue = evaluateNode(v.X, num) != evaluateNode(v.Y, num)
		case ast.Greater:
			conditionTrue = evaluateNode(v.X, num) > evaluateNode(v.Y, num)
		case ast.GreaterOrEqual:
			conditionTrue = evaluateNode(v.X, num) >= evaluateNode(v.Y, num)
		case ast.Less:
			conditionTrue = evaluateNode(v.X, num) < evaluateNode(v.Y, num)
		case ast.LessOrEqual:
			conditionTrue = evaluateNode(v.X, num) <= evaluateNode(v.Y, num)
		case ast.Reminder:
			rightVal := evaluateNode(v.Y, num)
			if rightVal == 0 {
				return 0
			}
			return evaluateNode(v.X, num) % rightVal
		}
	default:
		return 0
	}

	if conditionTrue {
		return 1
	}

	return 0
}
