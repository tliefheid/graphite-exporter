package main

import (
	"fmt"
	"os"
	"strings"

	logging "github.com/op/go-logging"
)

func trimAndReplace(s string) string {
	s = strings.Trim(s, " ")
	old := s
	s = strings.Replace(s, " ", "_", -1)
	new := s
	if new != old {
		fmt.Printf("WARN: parsed invalid config '%v' to valid '%v'. Change this in your config file!\n", old, new)
	}
	return s
}
func trimAndReplaceRef(s *string) {
	*s = trimAndReplace(*s)
	// *s = strings.Trim(*s, " ")
	// *s = strings.Replace(*s, " ", "_", -1)
}

func getKeyValue(input string, sep string) (string, string) {
	s := strings.Split(input, sep)
	key := trimAndReplace(s[0])
	val := trimAndReplace(s[1])
	return key, val
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func getLogLevel(input string) logging.Level {
	if input == "" {
		return logging.NOTICE
	}

	if strings.EqualFold(input, "debug") {
		return logging.DEBUG
	} else if strings.EqualFold(input, "info") {
		return logging.INFO
	} else if strings.EqualFold(input, "notice") {
		return logging.NOTICE
	} else if strings.EqualFold(input, "warn") {
		return logging.WARNING
	} else if strings.EqualFold(input, "error") {
		return logging.ERROR
	} else if strings.EqualFold(input, "critical") {
		return logging.CRITICAL
	}
	return logging.NOTICE
}
