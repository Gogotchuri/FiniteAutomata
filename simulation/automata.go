package simulation

import (
	"fmt"
	"github.com/gogotchuri/FiniteAutomata/parser"
)

type Transition map[string]*StateSet

type FiniteAutomata struct {
	States          []*State //Sort by State.ID
	Transitions     map[*State]Transition
	backwardConnect map[*State]Transition
	InitialState    *State
	AcceptingStates []*State
}

type State struct {
	ID          uint
	Name        string
	IsAccepting bool
}

func createAutomataForLiteral(literal rune) *FiniteAutomata {
	initialState := &State{ID: 0, Name: "q0", IsAccepting: false}
	acceptingState := &State{ID: 1, Name: "q1", IsAccepting: true}
	states := []*State{initialState, acceptingState}

	transitions := make(map[*State]Transition)
	backwardTransitions := make(map[*State]Transition)

	transitions[initialState] = make(Transition)
	backwardTransitions[acceptingState] = make(Transition)

	transitions[initialState][string(literal)] = CreateStateSet(acceptingState)
	backwardTransitions[acceptingState][string(literal)] = CreateStateSet(initialState)
	return &FiniteAutomata{
		States:          states,
		Transitions:     transitions,
		backwardConnect: backwardTransitions,
		InitialState:    initialState,
		AcceptingStates: []*State{acceptingState},
	}
}

func CreateFiniteAutomata(states []*State, transitions map[*State]Transition, initialState *State, acceptingStates []*State) *FiniteAutomata {
	return &FiniteAutomata{
		States:          states,
		Transitions:     transitions,
		InitialState:    initialState,
		AcceptingStates: acceptingStates,
	}
}

//Helpers
func (fa FiniteAutomata) String() string {
	transitionCount := fa.countTransitions()
	l1 := fmt.Sprintf("%d %d %d", len(fa.States), len(fa.AcceptingStates), transitionCount)
	l2 := fa.getAcceptingStates()
	l3 := fa.getTransitions()
	return fmt.Sprintf("%s\n%s\n%s", l1, l2, l3)
}

func (fa FiniteAutomata) countTransitions() int {
	transitionCount := 0
	for _, s := range fa.States {
		transitionCount += fa.countTransitionsForState(s)
	}
	return transitionCount
}

func (fa FiniteAutomata) countTransitionsForState(state *State) int {
	transitionCount := 0
	for _, t := range fa.Transitions[state] {
		transitionCount += t.Size()
	}
	return transitionCount
}

func (fa FiniteAutomata) getAcceptingStates() string {
	s := ""
	for i, a := range fa.AcceptingStates {
		if i == len(fa.AcceptingStates)-1 {
			s += fmt.Sprintf("%d", a.ID)
		} else {
			s += fmt.Sprintf("%d ", a.ID)
		}
	}
	return s
}

func (fa FiniteAutomata) getTransitions() string {
	res := ""
	for _, state := range fa.States {
		tc := fa.countTransitionsForState(state)
		if tc == 0 {
			res += fmt.Sprintf("%d\n", tc)
			continue
		}
		toAppend := fmt.Sprintf("%d ", tc)
		for symbol, toStates := range fa.Transitions[state] {
			for s2 := range toStates.Elements() {
				toAppend += fmt.Sprintf("%s %d ", symbol, s2.ID)
			}
		}
		res += toAppend[:len(toAppend)-1]
		res += "\n"
	}

	return res
}

func (fa FiniteAutomata) hasEpsilonTransitions() bool {
	for _, state := range fa.States {
		if fa.Transitions[state][parser.Epsilon] != nil && fa.Transitions[state][parser.Epsilon].Size() > 0 {
			return true
		}
	}
	return false
}
