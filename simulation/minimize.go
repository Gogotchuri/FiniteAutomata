package simulation

import (
	"fmt"
	"github.com/gogotchuri/FiniteAutomata/parser"
)

func (fa *FiniteAutomata) Minimize() *FiniteAutomata {
	newFA := fa.removeEpsilonTransitions()
	// Remove any remaining epsilon transitions
	for newFA.hasEpsilonTransitions() {
		newFA = newFA.removeEpsilonTransitions()
	}
	newFA.relabelStates()
	newFA = newFA.removeUnreachableStates()
	newFA.relabelStates()
	newFA.removeDeadStates()
	newFA.relabelStates()

	return newFA
}

func (fa *FiniteAutomata) removeEpsilonTransitions() *FiniteAutomata {
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

			if fa.countBackwardConnectionsForState(toState) == 1 &&
				fa.backwardConnect[toState][parser.Epsilon].Size() == 1 {
				fa.RemoveState(toState)
				if toState == fa.InitialState {
					fa.InitialState = s
				}
			}
		}
	}

	return fa
}

func (fa *FiniteAutomata) countBackwardConnectionsForState(state *State) int {
	count := 0
	for _, t := range fa.backwardConnect[state] {
		count += t.Size()
	}

	return count
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
				// Add backward transition from s to dst
				if _, ok := fa.backwardConnect[s]; !ok {
					fa.backwardConnect[s] = make(Transition)
				}
				if _, ok := fa.backwardConnect[s][symbol]; !ok {
					fa.backwardConnect[s][symbol] = CreateStateSet(dst)
				} else {
					fa.backwardConnect[s][symbol].Add(dst)
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
		}
	}
	if addsAccepted {
		fa.considerEpsilonsToAcceptStates()
	}
}

func (fa *FiniteAutomata) RemoveState(state *State) *FiniteAutomata {
	var index int
	for i, s := range fa.States {
		if s == state {
			index = i
			continue
		}
		if _, ok := fa.Transitions[s]; !ok {
			continue
		}
		for _, toStates := range fa.Transitions[s] {
			toStates.Remove(state)
		}
		if _, ok := fa.backwardConnect[s]; !ok {
			continue
		}
		for _, fromStates := range fa.backwardConnect[s] {
			fromStates.Remove(state)
		}
	}
	if index == len(fa.States)-1 {
		fa.States = fa.States[:index]
	} else {
		fa.States = append(fa.States[:index], fa.States[index+1:]...)
	}
	//Remove from accept states
	for i, s := range fa.AcceptingStates {
		if s == state {
			if index == len(fa.AcceptingStates)-1 {
				fa.AcceptingStates = fa.AcceptingStates[:i]
			} else {
				fa.AcceptingStates = append(fa.AcceptingStates[:i], fa.AcceptingStates[i+1:]...)
			}
			break
		}
	}
	delete(fa.Transitions, state)
	delete(fa.backwardConnect, state)
	return fa
}

func (fa *FiniteAutomata) relabelStates() {
	initialIndex := 0
	for i, s := range fa.States {
		if s == fa.InitialState {
			initialIndex = i
		}
	}
	fa.States[0], fa.States[initialIndex] = fa.States[initialIndex], fa.States[0]
	for i, s := range fa.States {
		s.ID = uint(i)
		s.Name = fmt.Sprintf("q%d", i)
	}
}

func (fa FiniteAutomata) removeUnreachableStates() *FiniteAutomata {
	// Remove unreachable states
	// by removing all states that are not reachable from the initial state
	// by performing a depth first search
	visited := CreateEmptyStateSet()
	fa.dfs(fa.InitialState, visited)
	newStates := make([]*State, visited.Size())
	i := 0
	for _, s := range fa.States {
		if visited.Contains(s) {
			newStates[i] = s
			i++
		}
	}
	newTransitions := make(map[*State]Transition)
	backwardTransitions := make(map[*State]Transition)
	for _, s := range newStates {
		if _, ok := fa.Transitions[s]; ok {
			newTransitions[s] = fa.Transitions[s]
		}
		if _, ok := fa.backwardConnect[s]; ok {
			backwardTransitions[s] = fa.backwardConnect[s]
		}
	}
	newFA := CreateFiniteAutomata(newStates, newTransitions, fa.InitialState, fa.AcceptingStates)
	newFA.backwardConnect = backwardTransitions
	return newFA
}

func (fa *FiniteAutomata) removeDeadStates() *FiniteAutomata {
	// Remove dead states
	// by removing all states that can never lead to an accepting state
	deadStates := CreateEmptyStateSet()
	for _, s := range fa.States {
		if !fa.canReachAcceptingState(s) {
			deadStates.Add(s)
		}
	}
	for s := range deadStates.Elements() {
		fa.RemoveState(s)
	}
	return fa
}
