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
