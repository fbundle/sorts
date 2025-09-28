package slices_util

func Map[T1 any, T2 any](input []T1, f func(T1) T2) []T2 {
	output := make([]T2, len(input))
	for i, v := range input {
		output[i] = f(v)
	}
	return output
}

func IndexedMap[T1 any, T2 any](input []T1, f func(int, T1) T2) []T2 {
	output := make([]T2, len(input))
	for i, v := range input {
		output[i] = f(i, v)
	}
	return output
}

func Reduce[T1 any, T2 any](input []T1, init T2, f func(T2, T1) T2) T2 {
	output := init
	for _, v := range input {
		output = f(output, v)
	}
	return output
}

func Reverse[T any](input []T) []T {
	output := make([]T, len(input))
	for i, v := range input {
		output[len(input)-1-i] = v
	}
	return output
}

func ForEach[T any](input []T, f func(T)) {
	for _, v := range input {
		f(v)
	}
}

func IndexedForEach[T any](input []T, f func(int, T)) {
	for i, v := range input {
		f(i, v)
	}
}
