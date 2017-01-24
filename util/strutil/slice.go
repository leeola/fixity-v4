package strutil

func DelFromSlice(l []string, s string) []string {
	for i, v := range l {
		if v == s {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func AddUnique(l []string, s string) []string {
	for _, v := range l {
		if v == s {
			return l
		}
	}
	return append(l, s)
}
