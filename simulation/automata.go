package simulation

import (
	"fmt"
	"github.com/gogotchuri/FiniteAutomata/parser"
)

type Transition map[string]*StateSet

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
	transitions[initialState][string(literal)] = CreateStateSet(acceptingState)
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
	for i, state := range fa.States {
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
			for s2 := range toStates.Elements() {
				if transitions[s][symbol] == nil {
					transitions[s][symbol] = CreateEmptyStateSet()
				}
				transitions[s][symbol].Add(s2)
			}
		}
	}
	for _, s := range fa2.States {
		for symbol, toStates := range fa2.Transitions[s] {
			for s2 := range toStates.Elements() {
				if transitions[s][symbol] == nil {
					transitions[s][symbol] = CreateEmptyStateSet()
				}
				transitions[s][symbol].Add(s2)
			}
		}
	}
	// Epsilon transitions from fa accept states to fa2 initial state
	for _, s := range fa.AcceptingStates {
		if _, ok := transitions[s]; !ok {
			transitions[s] = make(Transition)
		}
		if transitions[s][parser.Epsilon] == nil {
			transitions[s][parser.Epsilon] = CreateEmptyStateSet()
		}
		transitions[s][parser.Epsilon].Add(fa2.InitialState)
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
			for s2 := range toStates.Elements() {
				if transitions[s][symbol] == nil {
					transitions[s][symbol] = CreateEmptyStateSet()
				}
				transitions[s][symbol].Add(s2)
			}
		}
	}
	for _, s := range a2.States {
		for symbol, toStates := range a2.Transitions[s] {
			for s2 := range toStates.Elements() {
				if transitions[s][symbol] == nil {
					transitions[s][symbol] = CreateEmptyStateSet()
				}
				transitions[s][symbol].Add(s2)
			}
		}
	}
	if transitions[fa.InitialState][parser.Epsilon] == nil {
		transitions[fa.InitialState][parser.Epsilon] = CreateEmptyStateSet()
	}
	//Add epsilon transitions from fa initial state to a2 initial state
	transitions[fa.InitialState][parser.Epsilon].Add(a2.InitialState)
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
		if fa.Transitions[s][parser.Epsilon] == nil {
			fa.Transitions[s][parser.Epsilon] = CreateEmptyStateSet()
		}
		fa.Transitions[s][parser.Epsilon].Add(fa.InitialState)
	}
	acceptingStates = append(acceptingStates, fa.InitialState)
	return CreateFiniteAutomata(states, fa.Transitions, fa.InitialState, acceptingStates)
}

func (fa *FiniteAutomata) Minimize() *FiniteAutomata {
	newFA := fa.RemoveEpsilonTransitions()
	//newFA = newFA.RemoveUnreachableStates()
	//TODO implement
	return newFA
}

func (fa *FiniteAutomata) RemoveEpsilonTransitions() *FiniteAutomata {
	fa.considerEpsilonsToAcceptStates()
	// Remove epsilon transitions from to state A to state B
	// by adding every transition from state B to every other state to state A
	for _, s := range fa.States {
		if fa.Transitions[s] == nil {
			continue
		}
		if _, ok := fa.Transitions[s][parser.Epsilon]; !ok {
			fa.Transitions[s][parser.Epsilon] = CreateEmptyStateSet()
			continue
		}
		if fa.Transitions[s][parser.Epsilon] == nil {
			continue
		}
		for toState := range fa.Transitions[s][parser.Epsilon].Elements() {
			fa.Transitions[s][parser.Epsilon].Remove(toState)
			fa.appendTransition(s, toState, CreateEmptyStateSet())
		}
	}

	return fa
}

func (fa *FiniteAutomata) appendTransition(dst *State, src *State, visited *StateSet) {
	if _, ok := fa.Transitions[dst]; !ok {
		fa.Transitions[dst] = make(Transition)
	}
	// We have already visited this state
	if visited.Contains(src) {
		return
	}
	visited.Add(src)
	newEpsilonTransitions := CreateEmptyStateSet()
	for symbol, toStates := range fa.Transitions[src] {
		for s := range toStates.Elements() {
			if symbol == parser.Epsilon && s == dst {
				continue
			}
			if symbol == parser.Epsilon {
				// We will add further epsilon transitions to this state, so we don't need to add it to the transitions
				newEpsilonTransitions.Add(s)
			} else {
				if _, ok := fa.Transitions[dst][symbol]; !ok {
					fa.Transitions[dst][symbol] = CreateStateSet(s)
				} else {
					fa.Transitions[dst][symbol].Add(s)
				}
			}
		}
	}
	for src = range newEpsilonTransitions.Elements() {
		fa.appendTransition(dst, src, visited)
	}
}

func (fa *FiniteAutomata) considerEpsilonsToAcceptStates() {
	//epsilon transitions to accept states,
	//by making states accepting if they have epsilon transitions to accept states
	addsAccepted := false
	for _, s := range fa.States {
		if _, ok := fa.Transitions[s]; !ok {
			continue
		}
		if _, ok := fa.Transitions[s][parser.Epsilon]; !ok {
			continue
		}
		for s2 := range fa.Transitions[s][parser.Epsilon].Elements() {
			if !s2.IsAccepting || s.IsAccepting {
				continue
			}
			addsAccepted = true
			s.IsAccepting = true
			fa.AcceptingStates = append(fa.AcceptingStates, s)
			// Remove epsilon transition from s to s2
			//fa.Transitions[s][parser.Epsilon] = append(fa.Transitions[s][parser.Epsilon][:i], fa.Transitions[s][parser.Epsilon][i+1:]...)
		}
	}
	if addsAccepted {
		fa.considerEpsilonsToAcceptStates()
	}
}

/*
func (fa FiniteAutomata) RemoveUnreachableStates() *FiniteAutomata {
	// Remove unreachable states
	// by removing all states that are not reachable from the initial state
	// by performing a depth first search
	visited := make(map[*State]struct{})
	dfs(fa.InitialState, &visited)
	newStates := make([]*State, 0)
	for _, s := range fa.States {
		if _, ok := visited[s]; ok {
			newStates = append(newStates, s)
		}
	}
	newTransitions := make(Transition)
	for _, s := range newStates {
		newTransitions[s] = make(map[parser.Symbol][]*State)
		for symbol, toStates := range fa.Transitions[s] {
			newTransitions[s][symbol] = make([]*State, 0)
			for _, s2 := range toStates {
				if _, ok := visited[s2]; ok {
					newTransitions[s][symbol] = append(newTransitions[s][symbol], s2)
				}
			}
		}
	}
	return CreateFiniteAutomata(newStates, newTransitions, fa.InitialState, fa.AcceptingStates)
}
*/
