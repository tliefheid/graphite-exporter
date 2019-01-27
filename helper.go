package main

import "strings"

func trimAndReplace(s string) string {
	s = strings.Trim(s, " ")
	s = strings.Replace(s, " ", "_", -1)
	return s
}

func getKeyValue(input string, sep string) (string, string) {
	s := strings.Split(input, sep)
	key := trimAndReplace(s[0])
	val := trimAndReplace(s[1])
	return key, val
}
