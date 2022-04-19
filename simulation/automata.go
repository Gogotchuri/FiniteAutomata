package simulation

import (
	"FiniteAutomata/parser"
	"fmt"
)

type Transition map[string][]*State

type FiniteAutomata struct {
	States          []*State //Sort by State.ID
	Transitions     map[*State]Transition
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
	transitions[initialState] = make(Transition)
	transitions[initialState][string(literal)] = []*State{acceptingState}
	return &FiniteAutomata{
		States:          states,
		Transitions:     transitions,
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
		for _, t := range fa.Transitions[s] {
			transitionCount += len(t)
		}
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
	for i, state := range fa.States {
		toAppend := fmt.Sprintf("%d ", state.ID)
		for symbol, toStates := range fa.Transitions[state] {
			for _, s2 := range toStates {
				toAppend += fmt.Sprintf("%s %d ", symbol, s2.ID)
			}
		}
		if len(toAppend) < 3 {
			continue
		}
		if len(toAppend) > 1 {
			res += toAppend[:len(toAppend)-1]
		}
		if i != len(fa.States)-1 {
			res += "\n"
		}
	}

	return res
}

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
	for _, s := range states {
		transitions[s] = make(Transition)
	}
	for _, s := range fa.States {
		for symbol, toStates := range fa.Transitions[s] {
			for _, s2 := range toStates {
				transitions[s][symbol] = append(transitions[s][symbol], s2)
			}
		}
	}
	for _, s := range fa2.States {
		for symbol, toStates := range fa2.Transitions[s] {
			for _, s2 := range toStates {
				transitions[s][symbol] = append(transitions[s][symbol], s2)
			}
		}
	}
	// Epsilon transitions from fa accept states to fa2 initial state
	for _, s := range fa.AcceptingStates {
		if _, ok := transitions[s]; !ok {
			transitions[s] = make(Transition)
		}
		transitions[s][parser.Epsilon] = append(transitions[s][parser.Epsilon], fa2.InitialState)
	}
	return CreateFiniteAutomata(states, transitions, initialState, acceptingStates)
}

func (fa *FiniteAutomata) Union(a2 *FiniteAutomata) *FiniteAutomata {
	states := append(fa.States, a2.States...)
	for i, s := range states {
		s.ID = uint(i)
		s.Name = fmt.Sprintf("q%d", i)
	}
	transitions := make(map[*State]Transition)
	for _, s := range states {
		transitions[s] = make(Transition)
	}
	for _, s := range fa.States {
		for symbol, toStates := range fa.Transitions[s] {
			for _, s2 := range toStates {
				transitions[s][symbol] = append(transitions[s][symbol], s2)
			}
		}
	}
	for _, s := range a2.States {
		for symbol, toStates := range a2.Transitions[s] {
			for _, s2 := range toStates {
				transitions[s][symbol] = append(transitions[s][symbol], s2)
			}
		}
	}
	//Add epsilon transitions from fa initial state to a2 initial state
	transitions[fa.InitialState][parser.Epsilon] = append(transitions[fa.InitialState][parser.Epsilon], a2.InitialState)
	return CreateFiniteAutomata(states, transitions, fa.InitialState, append(fa.AcceptingStates, a2.AcceptingStates...))
}

func (fa *FiniteAutomata) Star() *FiniteAutomata {
	states := append(fa.States)
	fa.InitialState.IsAccepting = true
	acceptingStates := append(fa.AcceptingStates)
	for _, s := range acceptingStates {
		if _, ok := fa.Transitions[s]; !ok {
			fa.Transitions[s] = make(Transition)
		}
		fa.Transitions[s][parser.Epsilon] = append(fa.Transitions[s][parser.Epsilon], fa.InitialState)
	}
	acceptingStates = append(acceptingStates, fa.InitialState)
	return CreateFiniteAutomata(states, fa.Transitions, fa.InitialState, acceptingStates)
}

func (fa *FiniteAutomata) RemoveEpsilonTransitions() *FiniteAutomata {
	//TODO implement
	return fa
}