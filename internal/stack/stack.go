package stack

type Stack[T any] struct {
	arr []T
	c   int
}

func (s *Stack[T]) Push(v T) {
	if s.c == len(s.arr) {
		s.arr = append(s.arr, v)
	} else {
		s.arr[s.c] = v
	}

	s.c++
}

func (s *Stack[T]) Pop() T {
	s.c--
	return s.arr[s.c]
}

func (s *Stack[T]) Peek() T {
	return s.arr[s.c-1]
}

func (s *Stack[T]) Size() int {
	return s.c
}

func (s *Stack[T]) HasWith(f func(v T) bool) bool {
	for i := 0; i < s.c; i++ {
		if f(s.arr[i]) {
			return true
		}
	}

	return false
}

func New[T any](vs ...T) *Stack[T] {
	return &Stack[T]{
		arr: vs,
	}
}
