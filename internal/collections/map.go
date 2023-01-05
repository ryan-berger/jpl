package collections

func Map[T, U any](arr []T, fn func(T) U) []U {
	newArr := make([]U, len(arr))
	for i, t := range arr {
		newArr[i] = fn(t)
	}

	return newArr
}
