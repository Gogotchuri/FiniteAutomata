package simulation

import (
	"bufio"
	"fmt"
	"github.com/gogotchuri/FiniteAutomata/parser"
	"io"
	"strings"
)

func BuildAutomataFromDescriptionSTDIn(r io.Reader) (*FiniteAutomata, error) {
	var n, a, t int
	if _, err := fmt.Fscanf(r, "%d %d %d\n", &n, &a, &t); err != nil {
		return nil, err
	}
	states := make([]*State, n)
	transitions := make(map[*State]Transition, n)
	for i := 0; i < n; i++ {
		states[i] = &State{
			ID:          uint(i),
			Name:        fmt.Sprintf("q%d", i),
			IsAccepting: false,
		}
		transitions[states[i]] = make(Transition)
	}
	acceptStates := make([]*State, a)
	for i := 0; i < a; i++ {
		var s int
		if _, err := fmt.Fscanf(r, "%d", &s); err != nil {
			return nil, err
		}
		states[s].IsAccepting = true
		acceptStates[i] = states[s]
	}
	//Read transitions, so we have unified transition format
	transitionStrings := make([]string, n)
	inputReader := bufio.NewReader(r)
	for i := 0; i < n; i++ {
		input, err := inputReader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		input = strings.Trim(input, " \n")
		transitionStrings[i] = input
	}
	//Parse transitions
	for i, line := range transitionStrings {
		tokens := strings.Split(line, " ")
		var tc int
		if _, err := fmt.Sscanf(tokens[0], "%d", &tc); err != nil {
			return nil, err
		}
		if tc == 0 {
			continue
		}
		//Cut off the first token
		tokens = tokens[1:]
		for j := 0; j < tc; j++ {
			var toState int
			var sym string
			sym = tokens[j*2]
			if _, err := fmt.Sscanf(tokens[j*2+1], "%d", &toState); err != nil {
				return nil, err
			}
			if _, ok := transitions[states[i]][sym]; !ok {
				transitions[states[i]][sym] = CreateEmptyStateSet()
			}
			transitions[states[i]][sym].Add(states[toState])
		}
	}
	return CreateFiniteAutomata(states, transitions, states[0], acceptStates), nil
}

func BuildAutomataFromRegex(regex string) *FiniteAutomata {
	parseTree := parser.NewParseTree(regex)
	if parseTree == nil {
		return nil
	}
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
