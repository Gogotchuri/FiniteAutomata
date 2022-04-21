package simulation

type void struct{}
type Ss map[*State]void

type StateSet struct {
	m Ss
}

func CreateStateSet(elements ...*State) *StateSet {
	set := make(Ss)
	for _, element := range elements {
		set[element] = void{}
	}
	return &StateSet{set}
}

func CreateEmptyStateSet() *StateSet {
	set := make(Ss)
	return &StateSet{set}
}

func (s *StateSet) Add(state *State) {
	if s.m == nil {
		s.m = make(Ss)
	}
	s.m[state] = void{}
}

func (s *StateSet) Remove(state *State) {
	if s == nil || s.m == nil {
		return
	}
	delete(s.m, state)
}

func (s StateSet) Contains(state *State) bool {
	if s.m == nil {
		return false
	}
	_, ok := s.m[state]
	return ok
}

func (s StateSet) Size() int {
	if s.m == nil {
		return 0
	}
	return len(s.m)
}

func (s StateSet) IsEmpty() bool {
	return s.Size() == 0
}

func (s StateSet) Elements() Ss {
	if s.m == nil {
		s.m = make(Ss)
	}
	return s.m
}
