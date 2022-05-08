package plural

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/vorlif/spreak/internal/util"
)

type Forms struct {
	Nplurals int
	node     *node
}

func (f *Forms) IndexForN(n interface{}) int {
	if f.node == nil {
		return 0
	}

	num, err := util.ToNumber(n)
	if err != nil {
		return 0
	}

	result := f.node.evaluate(int(math.Round(num)))
	return result
}

type parser struct {
	s           *scanner
	lastToken   token  // last read token
	lastLiteral string // last read literal
	n           int    // buffer size (max=1)
}

func MustParse(rule string) *Forms {
	f, err := Parse(rule)
	if err != nil {
		panic(err)
	}

	return f
}

func Parse(rule string) (*Forms, error) {
	p := &parser{
		s: newScanner(strings.NewReader(rule)),
	}
	return p.Parse()
}

func (p *parser) Parse() (*Forms, error) {
	forms := &Forms{}

	if tok, lit := p.scanIgnoreWhitespace(); tok != nPlurals {
		return nil, fmt.Errorf("po parser: found %q, expected 'nplurals'", lit)
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != assign {
		return nil, fmt.Errorf("po parser: found %q, expected '=' after nplurals", lit)
	}

	//nolint:revive
	if tok, lit := p.scanIgnoreWhitespace(); tok != number {
		return nil, fmt.Errorf("po parser: found %q, expected '=' after nplurals", lit)
	} else {
		forms.Nplurals, _ = strconv.Atoi(lit)
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != semicolon {
		return nil, fmt.Errorf("po parser: found %q, expected ';' after 'nplurals=%d'", lit, forms.Nplurals)
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != plural {
		return nil, fmt.Errorf("po parser: found %q, expected 'plural' after 'nplurals=%d; '", lit, forms.Nplurals)
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != assign {
		return nil, fmt.Errorf("po parser: found %q, expected '=' after 'nplurals=%d; plural'", lit, forms.Nplurals)
	}

	if errScan := p.scanNext(); errScan != nil {
		return nil, errScan
	}

	n, err := p.expression()
	if err != nil {
		return nil, err
	}

	forms.node = n

	if p.lastToken != semicolon {
		return nil, fmt.Errorf("po parser: found %q, expected ';' at end", p.lastLiteral)
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != eof {
		return nil, fmt.Errorf("po parser: found %q, expected end", lit)
	}

	return forms, nil

}

func (p *parser) expression() (*node, error) {
	n, err := p.logicalOrExpression()
	if err != nil {
		return nil, err
	}

	if p.lastToken == question {
		questionNode := newNode(p.lastToken)

		if errScan := p.scanNext(); errScan != nil {
			return nil, errScan
		}

		condTrue, errCt := p.expression()
		if errCt != nil {
			return nil, errCt
		}
		questionNode.setNode(1, condTrue)

		if p.lastToken != colon {
			return nil, fmt.Errorf("po parser: found %q, expected \":\"", p.lastLiteral)
		}

		if errScan := p.scanNext(); errScan != nil {
			return nil, errScan
		}

		condFalse, errCf := p.expression()
		if errCf != nil {
			return nil, errCf
		}

		questionNode.setNode(2, condFalse)
		questionNode.setNode(0, n)
		return questionNode, nil
	}

	return n, nil
}

func (p *parser) logicalOrExpression() (*node, error) {
	ln, err := p.logicalAndExpression()
	if err != nil {
		return nil, err
	}

	if p.lastToken == logicalOr {
		orNode := newNode(p.lastToken)

		if errScan := p.scanNext(); errScan != nil {
			return nil, errScan
		}

		rn, errRn := p.logicalOrExpression() // right
		if errRn != nil {
			return nil, errRn
		}

		if rn.token == logicalOr {
			orNode.setNode(0, ln)
			orNode.setNode(1, rn.getNode(0))
			rn.setNode(0, orNode)
			return rn, nil
		}

		orNode.setNode(0, ln)
		orNode.setNode(1, rn)
		return orNode, nil
	}

	return ln, nil
}

func (p *parser) logicalAndExpression() (*node, error) {
	ln, err := p.equalityExpression() // left
	if err != nil {
		return nil, err
	}

	if p.lastToken == logicalAnd {
		andNode := newNode(p.lastToken) // up

		if errScan := p.scanNext(); errScan != nil {
			return nil, errScan
		}

		rn, errRn := p.logicalAndExpression()
		if errRn != nil {
			return nil, errRn
		}

		if rn.token == logicalAnd {
			andNode.setNode(0, ln)
			andNode.setNode(1, rn.getNode(0))
			rn.setNode(0, andNode)
			return rn, nil
		}

		andNode.setNode(0, ln)
		andNode.setNode(1, rn)
		return andNode, nil
	}

	return ln, nil
}

func (p *parser) equalityExpression() (*node, error) {
	n, err := p.relationalExpression()
	if err != nil {
		return nil, err
	}

	if p.lastToken == equal || p.lastToken == notEqual {
		compareNode := newNode(p.lastToken)

		if errScan := p.scanNext(); errScan != nil {
			return nil, errScan
		}

		re, errRe := p.relationalExpression()
		if errRe != nil {
			return nil, errRe
		}

		compareNode.setNode(1, re)
		compareNode.setNode(0, n)
		return compareNode, nil
	}

	return n, nil
}

func (p *parser) relationalExpression() (*node, error) {
	n, err := p.multiplicativeExpression()
	if err != nil {
		return nil, err
	}

	if p.lastToken == greater || p.lastToken == less || p.lastToken == greaterOrEqual || p.lastToken == lessOrEqual {
		compareNode := newNode(p.lastToken)

		if errScan := p.scanNext(); errScan != nil {
			return nil, errScan
		}

		me, errMe := p.multiplicativeExpression()
		if errMe != nil {
			return nil, errMe
		}

		compareNode.setNode(1, me)
		compareNode.setNode(0, n)
		return compareNode, nil
	}

	return n, nil
}

func (p *parser) multiplicativeExpression() (*node, error) {
	n, err := p.pmExpression()
	if err != nil {
		return nil, err
	}

	if p.lastToken == reminder {
		reminderNode := newNode(p.lastToken)

		if errScan := p.scanNext(); errScan != nil {
			return nil, errScan
		}

		pm, errPm := p.pmExpression()
		if errPm != nil {
			return nil, errPm
		}

		reminderNode.setNode(1, pm)
		reminderNode.setNode(0, n)
		return reminderNode, nil
	}

	return n, nil
}

func (p *parser) pmExpression() (*node, error) {

	if p.lastToken == variable {
		variableNode := newNode(p.lastToken)
		if errScan := p.scanNext(); errScan != nil {
			return nil, errScan
		}
		return variableNode, nil
	} else if p.lastToken == number {
		numberNode := newNode(p.lastToken)
		numberNode.value, _ = strconv.Atoi(p.lastLiteral)
		if errScan := p.scanNext(); errScan != nil {
			return nil, errScan
		}
		return numberNode, nil
	} else if p.lastToken == leftBracket {
		if errScan := p.scanNext(); errScan != nil {
			return nil, errScan
		}

		exprNode, errEp := p.expression()
		if errEp != nil {
			return nil, errEp
		}

		if p.lastToken != rightBracket {
			return nil, fmt.Errorf("found %q, expected )", p.lastLiteral)
		}

		if errScan := p.scanNext(); errScan != nil {
			return nil, errScan
		}

		return exprNode, nil
	} else {
		return nil, fmt.Errorf("found %q, expected something other", p.lastLiteral)
	}
}

func (p *parser) scanNext() error {
	if tok, _ := p.scanIgnoreWhitespace(); tok == eof || tok == failure {
		return errors.New("eof reached without result")
	}
	return nil
}

func (p *parser) scan() (tok token, lit string) {
	if p.n != 0 {
		p.n = 0
		return p.lastToken, p.lastLiteral
	}

	tok, lit = p.s.scan()

	p.lastToken, p.lastLiteral = tok, lit

	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *parser) scanIgnoreWhitespace() (tok token, lit string) {
	tok, lit = p.scan()
	if tok == whitespace {
		tok, lit = p.scan()
	}
	return
}
