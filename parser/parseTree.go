package parser

type ParseTree struct {
	Value rune
	Left  *ParseTree
	Right *ParseTree
}

func (pt ParseTree) IsLeaf() bool {
	return pt.Left == nil && pt.Right == nil
}

func (pt ParseTree) IsStar() bool {
	return pt.Value == '*'
}

func (pt ParseTree) IsConcat() bool {
	return pt.Value == '.'
}

func (pt ParseTree) IsUnion() bool {
	return pt.Value == '|'
}

func (pt *ParseTree) String() string {
	if pt == nil {
		return ""
	}
	value := string(pt.Value)
	if pt.Left == nil && pt.Right == nil {
		return string(pt.Value)
	}
	if pt.Left != nil {
		value = "/" + value
	}
	if pt.Right != nil {
		value = value + "\\"
	}
	return "(" + pt.Left.String() + value + pt.Right.String() + ")"
}
