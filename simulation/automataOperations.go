package simulation

import (
	"fmt"
	"github.com/gogotchuri/FiniteAutomata/parser"
)

func (fa *FiniteAutomata) Concat(fa2 *FiniteAutomata) *FiniteAutomata {
	states := append(fa.States, fa2.States...)
	for i, s := range states {
		s.ID = uint(i)
		s.Name = fmt.Sprintf("q%d", i)
		s.IsAccepting = false
	}
	initialState := fa.InitialState
	acceptingStates := fa2.AcceptingStates
	for _, s := range acceptingStates {
		s.IsAccepting = true
	}
	//Merge Transitions
	transitions := make(map[*State]Transition)
	backwardTransitions := make(map[*State]Transition)
	for _, s := range states {
		transitions[s] = make(Transition)
		backwardTransitions[s] = make(Transition)
	}
	for _, s := range fa.States {
		for symbol, toStates := range fa.Transitions[s] {
			for s2 := range toStates.Elements() {
				if transitions[s][symbol] == nil {
					transitions[s][symbol] = CreateEmptyStateSet()
				}
				if backwardTransitions[s2][symbol] == nil {
					backwardTransitions[s2][symbol] = CreateEmptyStateSet()
				}
				transitions[s][symbol].Add(s2)
				backwardTransitions[s2][symbol].Add(s)

			}
		}
	}
	for _, s := range fa2.States {
		for symbol, toStates := range fa2.Transitions[s] {
			for s2 := range toStates.Elements() {
				if transitions[s][symbol] == nil {
					transitions[s][symbol] = CreateEmptyStateSet()
				}
				if backwardTransitions[s2][symbol] == nil {
					backwardTransitions[s2][symbol] = CreateEmptyStateSet()
				}
				transitions[s][symbol].Add(s2)
				backwardTransitions[s2][symbol].Add(s)
			}
		}
	}
	// Epsilon transitions from fa accept states to fa2 initial state
	for _, s := range fa.AcceptingStates {
		if transitions[s][parser.Epsilon] == nil {
			transitions[s][parser.Epsilon] = CreateEmptyStateSet()
		}
		if backwardTransitions[fa2.InitialState][parser.Epsilon] == nil {
			backwardTransitions[fa2.InitialState][parser.Epsilon] = CreateEmptyStateSet()
		}
		transitions[s][parser.Epsilon].Add(fa2.InitialState)
		backwardTransitions[fa2.InitialState][parser.Epsilon].Add(s)
	}
	newFA := CreateFiniteAutomata(states, transitions, initialState, acceptingStates)
	newFA.backwardConnect = backwardTransitions
	return newFA
}

func (fa *FiniteAutomata) Union(a2 *FiniteAutomata) *FiniteAutomata {
	states := append(fa.States, a2.States...)
	for i, s := range states {
		s.ID = uint(i)
		s.Name = fmt.Sprintf("q%d", i)
	}
	transitions := make(map[*State]Transition)
	backwardTransitions := make(map[*State]Transition)
	for _, s := range states {
		transitions[s] = make(Transition)
		backwardTransitions[s] = make(Transition)
	}
	for _, s := range fa.States {
		for symbol, toStates := range fa.Transitions[s] {
			for s2 := range toStates.Elements() {
				if transitions[s][symbol] == nil {
					transitions[s][symbol] = CreateEmptyStateSet()
				}
				if backwardTransitions[s2][symbol] == nil {
					backwardTransitions[s2][symbol] = CreateEmptyStateSet()
				}
				transitions[s][symbol].Add(s2)
				backwardTransitions[s2][symbol].Add(s)
			}
		}
	}
	for _, s := range a2.States {
		for symbol, toStates := range a2.Transitions[s] {
			for s2 := range toStates.Elements() {
				if transitions[s][symbol] == nil {
					transitions[s][symbol] = CreateEmptyStateSet()
				}
				if backwardTransitions[s2][symbol] == nil {
					backwardTransitions[s2][symbol] = CreateEmptyStateSet()
				}
				transitions[s][symbol].Add(s2)
				backwardTransitions[s2][symbol].Add(s)
			}
		}
	}
	if transitions[fa.InitialState][parser.Epsilon] == nil {
		transitions[fa.InitialState][parser.Epsilon] = CreateEmptyStateSet()
	}
	if backwardTransitions[a2.InitialState][parser.Epsilon] == nil {
		backwardTransitions[a2.InitialState][parser.Epsilon] = CreateEmptyStateSet()
	}
	//Add epsilon transitions from fa initial state to a2 initial state
	transitions[fa.InitialState][parser.Epsilon].Add(a2.InitialState)
	backwardTransitions[a2.InitialState][parser.Epsilon].Add(fa.InitialState)
	newFA := CreateFiniteAutomata(states, transitions, fa.InitialState, append(fa.AcceptingStates, a2.AcceptingStates...))
	newFA.backwardConnect = backwardTransitions
	return newFA
}

func (fa *FiniteAutomata) Star() *FiniteAutomata {
	states := append(fa.States)
	fa.InitialState.IsAccepting = true
	acceptingStates := append(fa.AcceptingStates)
	for _, s := range acceptingStates {
		if _, ok := fa.Transitions[s]; !ok {
			fa.Transitions[s] = make(Transition)
		}

		if fa.Transitions[s][parser.Epsilon] == nil {
			fa.Transitions[s][parser.Epsilon] = CreateEmptyStateSet()
		}
		fa.Transitions[s][parser.Epsilon].Add(fa.InitialState)

		//Add backward transition from fa initial state to s
		if _, ok := fa.backwardConnect[fa.InitialState]; !ok {
			fa.backwardConnect[fa.InitialState] = make(Transition)
		}
		if fa.backwardConnect[fa.InitialState][parser.Epsilon] == nil {
			fa.backwardConnect[fa.InitialState][parser.Epsilon] = CreateEmptyStateSet()
		}
		fa.backwardConnect[fa.InitialState][parser.Epsilon].Add(s)
	}
	acceptingStates = append(acceptingStates, fa.InitialState)
	newFA := CreateFiniteAutomata(states, fa.Transitions, fa.InitialState, acceptingStates)
	newFA.backwardConnect = fa.backwardConnect
	return newFA
}
