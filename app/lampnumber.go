package app

import "strings"

func getLampNumber(file string) string {
	s := strings.TrimSuffix(file, ".txt")
	l := len(s)
	return s[l-5:]
}
