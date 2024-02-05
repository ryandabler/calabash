package slice

func Map[T any, V any](as []T, f func(T) (V, error)) ([]V, error) {
	bs := make([]V, len(as))

	for i, a := range as {
		b, err := f(a)

		if err != nil {
			return nil, err
		}

		bs[i] = b
	}

	return bs, nil
}

func Fold[T any, V any](as []T, i V, f func(T, V, int) V) V {
	for idx, a := range as {
		i = f(a, i, idx)
	}

	return i
}

func Last[T any](as []T) (T, bool) {
	var empty T

	if len(as) == 0 {
		return empty, false
	}

	return as[len(as)-1], true
}
