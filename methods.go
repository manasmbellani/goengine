package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/go-resty/resty"
)

// LinBashPath is the Linux shell basepath
const LinBashPath string = "/bin/bash"

// WinCmdPath is Windows command path is the Linux shell basepath
const WinCmdPath string = "cmd.exe"

// DefUserAgent - Default user agent string
const DefUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36"

// WebTimeout is the timeout for web requests
const WebTimeout = 5

// execMethod is generallly used to execute particular commands
func execMethod(target Target, checkID string, methodID string,
	method MethodStruct, outfolder string) {

	methodType := method.Type

	log.Printf("[*] Executing checkID: %s, methodID: %s of type: %s on target: %+v\n",
		checkID, methodID, methodType, target)
	if methodType == "cmd" {
		execCmd(target, checkID, methodID, method, outfolder)
	} else if methodType == "webrequest" {
		execWebRequest(target, checkID, methodID, method, outfolder)
	} else {
		log.Printf("[-] Unknown method: %s, %s, %s\n", checkID, methodID, methodType)
	}

}

// execCmd is used to execute shell commands and return the results
func execCmd(target Target, checkID string, methodID string,
	method MethodStruct, outfolder string) {

	// Read the necessary variables to execute
	cmdDir := method.CmdDir
	cmds := method.Cmds
	regex := method.Regex
	outfile := method.Outfile

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

	// If matching regex found, then print the result
	if shouldNotify(totalOut, regex) {
		fmt.Printf("[%s-%s] %s\n", checkID, methodID, target.Target)
	} else {
		if outfile != "" {
			writeToOutfile(outfile, outfolder, totalOut, target)
		}
	}

}

// execWebRequest is used to execute web requests on a specific target given the
// relevant method
func execWebRequest(target Target, checkID string, methodID string,
	method MethodStruct, outfolder string) {
	// Read vars for processing
	urls := method.Urls
	httpMethod := method.HTTPMethod
	regex := method.Regex
	mheaders := method.Headers
	mbody := method.Body
	outfile := method.Outfile

	totalOut := ""

	// Create the restyClient for making web requests in this thread
	restyClient := resty.New()
	restyClient.SetTimeout(time.Duration(WebTimeout) * time.Second)

	for _, urlToCheck := range urls {

		// Determine if HTTP method is supported
		httpMethod := strings.ToUpper(httpMethod)
		if httpMethod == "" {
			httpMethod = "GET"
		}

		// Currently, we only support specific methods
		if (httpMethod != "GET") && (httpMethod != "POST") {
			log.Printf("[-] Unsupported HTTP method: %s\n", httpMethod)
			break
		}

		// Build the URL to request + save it
		urlToCheckSub := subTargetParams(urlToCheck, target)

		// Set the headers and X-Forwarded-For/X-Forwarded-Host
		headers := make(map[string]string)
		headers["User-Agent"] = DefUserAgent
		headers["X-Forwarded-For"] = "127.0.0.1"
		headers["X-Forwarded-Host"] = "127.0.0.1"
		for _, h := range mheaders {
			headers[h.Name] = h.Value
		}
		restyClient.SetHeaders(headers)

		// Prepare POST body via provided names, values params
		body := make(map[string]string)
		if mbody != nil {
			for _, bodySet := range mbody {
				name := bodySet.Name
				value := bodySet.Value
				body[name] = value
			}
		}

		// Verbose message to be printed to let the user know
		log.Printf("[*] Make %s request to URL: %s\n", httpMethod,
			urlToCheckSub)
		var errResty error
		var respResty *resty.Response
		if httpMethod == "POST" {
			respResty, errResty = restyClient.R().SetBody(body).Post(urlToCheckSub)
		} else {
			respResty, errResty = restyClient.R().Get(urlToCheckSub)
		}

		// Check if there was an error
		if errResty != nil {
			log.Println("[-] Error making request to URL: ",
				urlToCheckSub, " Error: ", errResty)
		}
		log.Printf("[*] Getting the raw HTTP request")
		if errResty != nil {
			fmt.Println(errResty)
		}

		if respResty != nil {

			// Read the response body
			respBody := respResty.String()

			// Read the response status code as string
			statusCode := respResty.StatusCode()

			// Read the response headers as string
			respHeaders := respResty.Header()
			respHeadersStr := ""
			for k, v := range respHeaders {
				s := fmt.Sprintf("%s:%s", k, strings.Join(v, ","))
				respHeadersStr += s + ";"
			}

			// Combine status code, response headers and body
			requestOut := fmt.Sprintf("%d\n%s\n%s\n", statusCode,
				respHeadersStr, respBody)

			// If matching regex found, then print the result
			if shouldNotify(requestOut, regex) {
				fmt.Printf("[%s-%s] %s\n", checkID, methodID, urlToCheckSub)
			} else {
				if outfile != "" {
					writeToOutfile(outfile, outfolder, totalOut, target)
				}
			}

			totalOut += requestOut
		}
	}
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
