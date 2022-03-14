package util

func DeleteString(a []string, s string) ([]string) {
	for i, e := range a {
		if e == s {
			return append(a[:i], a[i+1:]...)
		}
	}
	return a
}

func AddString(a []string, s string) ([]string) {
	for _, e := range a {
		if e == s {
			return a
		}
	}
	return append(a, s)
}

func HasString(a []string, s string) (bool) {
	for _, e := range a {
		if e == s {
			return true
		}
	}
	return false
}

func MakeStringSet(a []string) (b []string) {
	aMap := make(map[string]bool)
	for _, e := range a {
		if _, exists := aMap[e]; !exists {
			aMap[e] = true
			b = append(b, e)
		}
	}

	return b
}

func DeleteInt(a []int64, s int64) ([]int64) {
	for i, e := range a {
		if e == s {
			return append(a[:i], a[i+1:]...)
		}
	}
	return a
}

func AddInt(a []int64, s int64) ([]int64) {
	for _, e := range a {
		if e == s {
			return a
		}
	}
	return append(a, s)
}

func HasInt(a []int64, s int64) (bool) {
	for _, e := range a {
		if e == s {
			return true
		}
	}
	return false
}
