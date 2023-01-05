package generator

func toArray[T comparable, V any](src map[T]V) (r []V) {
	r = make([]V, 0, len(src))
	for _, v := range src {
		r = append(r, v)
	}
	return r
}
