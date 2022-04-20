package simulation

import (
	"FiniteAutomata/parser"
)

func BuildAutomataFromDescription(description string) *FiniteAutomata {
	//tODO: implement
	return nil
}

func BuildAutomataFromRegex(regex string) *FiniteAutomata {
	parseTree := parser.NewParseTree(regex)
	return buildAutomataFromParseTree(parseTree).Minimize()
}

func buildAutomataFromParseTree(parseTree *parser.ParseTree) *FiniteAutomata {
	if parseTree.IsLeaf() {
		return createAutomataForLiteral(parseTree.Value)
	}
	if parseTree.IsConcat() {
		a1 := buildAutomataFromParseTree(parseTree.Left)
		a2 := buildAutomataFromParseTree(parseTree.Right)
		return a1.Concat(a2)
	}

	if parseTree.IsUnion() {
		a1 := buildAutomataFromParseTree(parseTree.Left)
		a2 := buildAutomataFromParseTree(parseTree.Right)
		return a1.Union(a2)
	}

	if parseTree.IsStar() {
		a1 := buildAutomataFromParseTree(parseTree.Left)
		r := a1.Star()
		return r
	}

	return nil
}
