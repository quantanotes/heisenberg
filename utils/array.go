package utils

func Equals[T comparable](as, bs []T) bool {
	if len(as) != len(bs) {
		return false
	}
	for i, a := range as {
		if a != bs[i] {
			return false
		}
	}
	return true
}

func Contains[T comparable](array []T, element T) bool {
	for _, e := range array {
		if e == element {
			return true
		}
	}
	return false
}

func Intersect[T comparable](as, bs []T) []T {
	intersect := make([]T, 0, len(bs))
	for _, a := range as {
		if Contains(bs, a) {
			intersect = append(intersect, a)
		}
	}
	return intersect
}

func Except[T comparable](as, bs []T) []T {
	except := make([]T, 0, len(bs))
	for _, a := range as {
		if !Contains(bs, a) {
			except = append(except, a)
		}
	}
	return except
}

func IntersectExcept[T comparable](as, bs []T) ([]T, []T) {
	intersect := make([]T, 0, len(bs))
	except := make([]T, 0, len(bs))
	for _, a := range as {
		if Contains(bs, a) {
			intersect = append(intersect, a)
		} else {
			except = append(except, a)
		}
	}
	return intersect, except
}

func RemoveAt[T any](array []T, index uint) []T {
	return append(array[:index], array[index+1:]...)
}

func RemoveWhere[T comparable](array []T, element T) []T {
	for i, e := range array {
		if e == element {
			return RemoveAt(array, uint(i))
		}
	}
	return array
}

func KeysAsArray[T comparable](m map[T]any) []T {
	ks := make([]T, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}
