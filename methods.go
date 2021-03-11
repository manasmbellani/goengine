package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// LinBashPath is the Linux shell basepath
const LinBashPath string = "/bin/bash"

// WinCmdPath is Windows command path is the Linux shell basepath
const WinCmdPath string = "cmd.exe"

// execMethod is generallly used to execute particular commands
func execMethod(target Target, checkID string, methodID string,
	method MethodStruct) {

	methodType := method.Type
	// outfile := method.Outfile
	regex := method.Regex

	// Store the output result of executing function
	out := ""

	log.Printf("[*] Executing checkID: %s, methodID: %s on target: %+v\n",
		checkID, methodID, target)
	if methodType == "cmd" {
		out = execCmd(target, method)
	} else if methodType == "webrequest" {
		log.Printf("Executing webrequest")
	} else if methodType == "" {
		// Do nothing
	} else {
		log.Fatalf("Unknown method: methodType")
	}

	if regex != "" {
		// If matching regex found, then print the result
		if shouldNotify(out, regex) {
			fmt.Printf("[%s-%s] %s\n", checkID, methodID, target.Target)
		}
	}

}

// execCmd is used to execute shell commands and return the results
func execCmd(target Target, method MethodStruct) string {
	// Read the necessary variables to execute
	cmdDir := method.CmdDir
	cmds := method.Cmds
	owd, _ := os.Getwd()
	if cmdDir == "" {
		cmdDir = owd
	}

	// Check if cmddir exists - otherwise, cannot execute anything
	if _, err := os.Stat(cmdDir); os.IsNotExist(err) {
		log.Printf("[-] Dir Path: %s not found", cmdDir)
	}

	// Build the commands
	joinedCmds := strings.Join(cmds, "; ")
	joinedCmds = fmt.Sprintf("cd %s; "+subTargetParams(joinedCmds, target)+"; cd %s",
		cmdDir, owd)

	// Let user know we are executing command
	log.Printf("[*] Executing command: %s\n", joinedCmds)

	// Determine the command to execute
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command(WinCmdPath, "/c", joinedCmds)
	default:
		cmd = exec.Command(LinBashPath, "-c", joinedCmds)
	}

	// Execute the command and get the output and error message
	out, err := cmd.CombinedOutput()
	var outStr, errStr, totalOut string
	if out == nil {
		outStr = ""
	} else {
		outStr = string(out)
	}

	if err == nil {
		errStr = ""
	} else {
		errStr = string(err.Error())
	}

	totalOut = (outStr + "\n" + errStr)
	return totalOut
}

// execNotes is used to print the notes given the target and the method details
func execNotes(target *Target, method *MethodStruct) {
	// Read the necessary variables to print notes
	notes := method.Notes
	log.Println("[!] Notes:")
	notesToPrint := strings.Split(notes, "\n")
	for _, note := range notesToPrint {
		noteToPrint := subTargetParams(note, *target)
		log.Println("[!] " + noteToPrint)
	}
}
