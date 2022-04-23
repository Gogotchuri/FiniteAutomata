package parser

import "strings"

func NewParseTree(regex string) *ParseTree {
	if regex == "" {
		return nil
	}
	regex = insertConcatenation(regex)
	return createRegexParseTree(regex)
}

//operator precedence (*, ., |)
//Recursive function to create the parse tree
func createRegexParseTree(regex string) *ParseTree {
	//Remove parenthesis if whole regex is enclosed in them
	if rx, ok := removeOuterParenthesis(regex); ok {
		return createRegexParseTree(rx)
	}
	if len(regex) == 0 {
		return nil
	}
	if len(regex) == 1 {
		return &ParseTree{Value: rune(regex[0])}
	}
	//Next operator to be evaluated
	beforeOperator, operator, afterOperator := getNextOperator(regex)

	pr := &ParseTree{}
	// We shouldn't encounter operator as a first character
	if isIn(rune(regex[0]), operators) {
		panic("Unexpected token encountered")
	}
	if !isIn(operator, operators) {
		panic(regex + " Invalid regex") //Expected operator after literal
	}

	pr.Value = operator
	pr.Left = createRegexParseTree(beforeOperator)
	if operator != '*' {
		pr.Right = createRegexParseTree(afterOperator)
	}
	return pr
}

func getNextOperator(regex string) (string, rune, string) {
	//Find the next operator to use in reverse order of precedence
	if ind := tryToFindOperator('|', regex); ind != -1 {
		return regex[:ind], '|', regex[ind+1:]
	} else if ind := tryToFindOperator('.', regex); ind != -1 {
		return regex[:ind], '.', regex[ind+1:]
	} else if ind := tryToFindOperator('*', regex); ind != -1 {
		return regex[:ind], '*', regex[ind+1:]
	}
	panic("Invalid regex")
}

func tryToFindOperator(op rune, regex string) int {
	depth := 0
	for i, curr := range regex {
		if curr == '(' {
			depth++
		} else if curr == ')' {
			depth--
		} else if curr == op && depth == 0 {
			return i
		}
	}
	return -1
}

func removeOuterParenthesis(regex string) (string, bool) {
	if len(regex) < 2 || regex[0] != '(' || regex[len(regex)-1] != ')' {
		return regex, false
	}
	if len(regex) == 3 {
		return regex[1 : len(regex)-1], true
	}
	depth := 0
	for i, curr := range regex {
		if curr == '(' {
			depth++
		} else if curr == ')' {
			depth--
		}
		//If until the last character, depth is 0, no outer parenthesis exists
		if depth == 0 && i != len(regex)-1 {
			return regex, false
		}
	}
	return regex[1 : len(regex)-1], true
}

// insertConcatenation Replaces all concatenations which are empty symbols with a single .
func insertConcatenation(regex string) string {
	regex = strings.ReplaceAll(regex, "()", Epsilon)
	var prev, curr rune
	res := string(regex[0])
	for i := 1; i < len(regex); i++ {
		prev = rune(regex[i-1])
		curr = rune(regex[i])
		if isIn(prev, literals+")*") && isIn(curr, literals+"(") {
			res += "." + string(curr)
		} else {
			res += string(curr)
		}
	}
	return res
}

func isIn(c rune, s string) bool {
	for _, r := range s {
		if c == r {
			return true
		}
	}
	return false
}
