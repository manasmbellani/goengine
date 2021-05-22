package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func eCmd(cmds []string, cmdDir string) string {
	owd, _ := os.Getwd()

	// Check if cmddir exists - otherwise, cannot execute anything
	if cmdDir != "" {
		if _, err := os.Stat(cmdDir); os.IsNotExist(err) {
			log.Fatalf("[-] Dir Path: %s not found", cmdDir)
		}
	}

	// Build the commands
	joinedCmds := strings.Join(cmds, ";")
	if cmdDir != "" {
		joinedCmds = fmt.Sprintf("cd %s; "+joinedCmds+"; cd %s",
			cmdDir, owd)
	}

	// Let user know we are executing command
	log.Printf("[*] Executing command: %s\n", joinedCmds)

	// Determine the command to execute
	var cmdObj *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmdObj = exec.Command(WinCmdPath, "/c", joinedCmds)
	default:
		cmdObj = exec.Command(LinBashPath, "-c", joinedCmds)
	}

	// Execute the command and get the output and error message
	out, err := cmdObj.CombinedOutput()
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
