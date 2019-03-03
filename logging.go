package main

import (
	"log"
)

func logMessage(format string, a ...interface{}) {
	if DebugLogging == true {
		log.Printf(format+"\n", a...)
	}
}
