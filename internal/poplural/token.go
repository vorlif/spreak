package poplural

type token string

const (
	eof            token = "eof"
	whitespace     token = "ws"
	failure        token = "failure"
	number         token = "number" // 1, 2, 100, etc.
	variable       token = "n"
	plural         token = "plural="
	nPlurals       token = "nplurals="
	equal          token = "=="
	assign         token = "="
	greater        token = ">"
	greaterOrEqual token = ">="
	less           token = "<"
	lessOrEqual    token = "<="
	reminder       token = "%"
	notEqual       token = "!="
	logicalAnd     token = "&&"
	logicalOr      token = "||"
	question       token = "?"
	colon          token = ":"
	semicolon      token = ";"
	leftBracket    token = "("
	rightBracket   token = ")"
)
