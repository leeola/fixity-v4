package strutil

func ShortHash(s string, l int) string {
	if s == "" {
		return ""
	}

	if l > len(s) {
		return s
	}

	return s[len(s)-l:]
}
