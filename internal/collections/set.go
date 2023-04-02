package collections

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(el T) {
	s[el] = struct{}{}
}

func (s Set[T]) Remove(el T) {
	delete(s, el)
}

func (s Set[T]) Contains(el T) bool {
	_, ok := s[el]
	return ok
}

func (s Set[T]) ToSlice() []T {
	sl := []T{}
	for el := range s {
		sl = append(sl, el)
	}
	return sl
}

func (s Set[T]) Clone() Set[T] {
	newSet := Set[T]{}
	for el := range s {
		newSet.Add(el)
	}
	return newSet
}

func (s Set[T]) Equals(rhs Set[T]) bool {
	if len(s) != len(rhs) {
		return false
	}
	rhsClone := rhs.Clone()
	for el := range s {
		rhsClone.Remove(el)
	}
	return len(rhsClone) == 0
}

func SliceToSet[T comparable](sl []T) Set[T] {
	set := Set[T]{}
	for _, el := range sl {
		set.Add(el)
	}
	return set
}
