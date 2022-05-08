package plural

type node struct {
	token token
	nodes [3]*node
	value int // Only if token is of type number
}

func newNode(t token) *node {
	return &node{
		token: t,
		nodes: [3]*node{},
	}
}

func (n *node) setNode(idx int, nd *node) {
	n.nodes[idx] = nd
}

func (n *node) getNode(idx int) *node {
	return n.nodes[idx]
}

func (n *node) evaluate(val int) int {
	var conditionTrue bool
	switch n.token {
	// leaf
	case number:
		return n.value
	case variable:
		return val
		// 2 args
	case equal:
		conditionTrue = n.nodes[0].evaluate(val) == n.nodes[1].evaluate(val)
	case notEqual:
		conditionTrue = n.nodes[0].evaluate(val) != n.nodes[1].evaluate(val)
	case greater:
		conditionTrue = n.nodes[0].evaluate(val) > n.nodes[1].evaluate(val)
	case greaterOrEqual:
		conditionTrue = n.nodes[0].evaluate(val) >= n.nodes[1].evaluate(val)
	case less:
		conditionTrue = n.nodes[0].evaluate(val) < n.nodes[1].evaluate(val)
	case lessOrEqual:
		conditionTrue = n.nodes[0].evaluate(val) <= n.nodes[1].evaluate(val)
	case logicalAnd:
		leftTrue := n.nodes[0].evaluate(val) == 1
		rightTrue := n.nodes[1].evaluate(val) == 1
		conditionTrue = leftTrue && rightTrue
	case logicalOr:
		conditionTrue = n.nodes[0].evaluate(val) == 1 || n.nodes[1].evaluate(val) == 1
	case reminder:
		rightVal := n.nodes[1].evaluate(val)
		if rightVal == 0 {
			return 0
		}
		return n.nodes[0].evaluate(val) % rightVal
	// 3 args
	case question:
		if n.nodes[0].evaluate(val) == 1 {
			return n.nodes[1].evaluate(val)
		}
		return n.nodes[2].evaluate(val)
	default:
		return 0
	}

	if conditionTrue {
		return 1
	}

	return 0
}
