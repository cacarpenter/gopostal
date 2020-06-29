// simple util functions until I find a library I'd rather use
package util

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
)

func ReplaceVariables(raw string, subVals map[string]string) string {
	r, err := regexp.Compile("\\{\\{\\w+\\}\\}")
	if err != nil {
		log.Panicln(err)
	}

	var buffer bytes.Buffer
	fasi := r.FindAllStringIndex(raw, -1)
	prevIdx := 0
	for _, matchIndices := range fasi {
		// TODO should be able to remove this
		if len(matchIndices) != 2 {
			panic(fmt.Errorf("Thought length would be 2 but it was %d\n", len(matchIndices)))
		}

		// chars before this match
		if matchIndices[0] > prevIdx {
			buffer.WriteString(raw[prevIdx:matchIndices[0]])
		}
		varstr := raw[matchIndices[0]+2 : matchIndices[1]-2]
		val := subVals[varstr]
		buffer.WriteString(val)

		prevIdx = matchIndices[1]
	}
	// some literal text remains
	if prevIdx < len(raw)-1 {
		buffer.WriteString(raw[prevIdx:len(raw)])
	}

	return buffer.String()
}

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

func Map2Array(m map[string]string) [][]string {
	arr := make([][]string, len(m))
	i := 0
	for k, v := range m {
		arr[i] = make([]string, 2)
		arr[i][0] = k
		arr[i][1] = v
		i++
	}
	return arr
}
