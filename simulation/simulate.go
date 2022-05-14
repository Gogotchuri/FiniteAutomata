package simulation

import (
	"fmt"
	"github.com/gogotchuri/FiniteAutomata/parser"
	"strings"
)

func (fa *FiniteAutomata) Simulate(input string) (bool, string) {
	input = strings.ReplaceAll(input, "()", parser.Epsilon)
	wasAccepted := make([]bool, len(input))
	if b := fa.run(input, fa.InitialState, wasAccepted); b && len(input) == 0 {
		return true, ""
	}
	resString := ""
	for _, b := range wasAccepted {
		if b {
			resString += "Y"
		} else {
			resString += "N"
		}
	}
	fmt.Println(resString)
	if len(resString) == 0 {
		return true, ""
	}
	return wasAccepted[len(wasAccepted)-1], resString
}

func (fa *FiniteAutomata) run(input string, currentState *State, acceptStates []bool) bool {
	if currentState == nil {
		currentState = fa.InitialState
	}
	curIndex := len(acceptStates) - len(input) - 1
	if currentState.IsAccepting && curIndex >= 0 {
		acceptStates[curIndex] = true
	}
	//Base case, if we reach a final state and string ends, return true
	if currentState.IsAccepting && len(input) == 0 {
		return true
	}
	var nextStates, epsilonStates *StateSet
	//Get next states, given the current state and the input
	if len(input) > 0 { //If input is empty we will just get epsilon transitions
		nextStates = fa.Transitions[currentState][string(input[0])]
	}
	epsilonStates = fa.Transitions[currentState][parser.Epsilon]
	//If there are no next states, return false
	if (nextStates == nil || nextStates.IsEmpty()) && (epsilonStates == nil || epsilonStates.IsEmpty()) {
		return false
	}
	if epsilonStates != nil {
		//Try to run the epsilon states
		for epsilonState := range epsilonStates.Elements() {
			if fa.run(input, epsilonState, acceptStates) {
				return true
			}
		}
	}
	if nextStates != nil {
		//Try to run the next states
		for nextState := range nextStates.Elements() {
			if fa.run(input[1:], nextState, acceptStates) {
				return true
			}
		}
	}
	//If we reach this point, we have failed to accept the input in the current state, return false
	return false
}
