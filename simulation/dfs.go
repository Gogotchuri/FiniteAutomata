package simulation

func (fa *FiniteAutomata) dfs(state *State, visited *StateSet) {
	visited.Add(state)
	for _, toStates := range fa.Transitions[state] {
		for nextState := range toStates.Elements() {
			if !visited.Contains(nextState) {
				fa.dfs(nextState, visited)
			}
		}
	}
}

func (fa *FiniteAutomata) canReachAcceptingState(state *State) bool {
	visited := CreateEmptyStateSet()
	fa.dfs(state, visited)
	for _, acceptingState := range fa.AcceptingStates {
		if visited.Contains(acceptingState) {
			return true
		}
	}
	return false
}
