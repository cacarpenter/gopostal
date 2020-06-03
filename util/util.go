// simple util functions until I find a library I'd rather use
package util

import "fmt"

func MaxLen(s1, s2 string) int {
	if len(s1) > len(s2) {
		return len(s1)
	}
	return len(s2)
}

func StringOf(c rune, l int) string {
	rs := make([]rune, l)
	for i := 0; i < l; i++ {
		rs[i] = c
	}
	return string(rs)
}

func Pad(s string, l int) string {
	return fmt.Sprintf("%-*s", l, s)
}
