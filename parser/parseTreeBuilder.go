package parser

import "strings"

func NewParseTree(regex string) *ParseTree {
	regex = InsertConcatenation(regex)
	return createRegexParseTree(regex)
}

//TODO operator precedence (*, ., |)
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
	if IsIn(rune(regex[0]), operators) {
		panic("Unexpected token encountered")
	}
	if !IsIn(operator, operators) {
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

func createParseTreeAfterParenthesis(regex string) *ParseTree {
	pr := &ParseTree{}
	findMatch := findMatchingParen(regex)
	if findMatch >= len(regex)-1 {
		pr = createRegexParseTree(regex[1 : len(regex)-1])
	} else {
		operator := rune(regex[findMatch+1])
		if operator == '*' {
			return createParseTreeForStar(regex[findMatch+2:], regex[1:findMatch])
		} else {
			pr.Value = rune(regex[findMatch+1])
			pr.Left = createRegexParseTree(regex[1:findMatch])
			pr.Right = createRegexParseTree(regex[findMatch+2:])
		}
	}
	return pr
}

func findMatchingParen(regex string) int {
	count := 1
	for i := 1; i < len(regex); i++ {
		if regex[i] == '(' {
			count++
		} else if regex[i] == ')' {
			count--
		}
		if count == 0 {
			return i
		}
	}
	panic("Unmatched parenthesis")
}

//Create a parse tree for a star (unary operator)
func createParseTreeForStar(regexAfterStar, leftSide string) *ParseTree {
	pr := &ParseTree{}
	pr.Value = '*'
	pr.Left = createRegexParseTree(leftSide)
	pr.Right = nil
	if len(regexAfterStar) == 0 {
		return pr
	}
	// a* type expression is taken already, we are expecting some kind of operator
	if !IsIn(rune(regexAfterStar[0]), operators) {
		panic("Invalid regex, expected operator after star") //Expected operator after literal
	}
	upper := &ParseTree{
		Value: rune(regexAfterStar[0]),
		Left:  pr,
		Right: createRegexParseTree(regexAfterStar[1:]),
	}
	return upper
}

// InsertConcatenation Replaces all concatenations which are empty symbols with a single .
func InsertConcatenation(regex string) string {
	regex = strings.ReplaceAll(regex, "()", Epsilon)
	var prev, curr rune
	res := string(regex[0])
	for i := 1; i < len(regex); i++ {
		prev = rune(regex[i-1])
		curr = rune(regex[i])
		if IsIn(prev, literals+")*") && IsIn(curr, literals+"(") {
			res += "." + string(curr)
		} else {
			res += string(curr)
		}
	}
	return res
}

func IsIn(c rune, s string) bool {
	for _, r := range s {
		if c == r {
			return true
		}
	}
	return false
}
