package main

import (
	"log"
)

// execMethod is generallly used to execute particular commands
func execMethod(target *Target, method *MethodStruct) {
	methodType := method.Type
	if methodType == "cmd" {
		log.Printf("Executing cmd")
	} else if methodType == "webrequest" {
		log.Printf("Executing webrequest")
	} else if methodType == "" {
		// Do nothing
	} else {
		log.Fatalf("Unknown method: methodType")
	}
}
