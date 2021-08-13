package main

import (
	"fmt"
	"log"
	"strings"
	"crypto/tls"
	"time"

	"github.com/go-resty/resty/v2"
)

// LinBashPath is the Linux shell basepath
const LinBashPath string = "/bin/bash"

// WinCmdPath is Windows command path is the Linux shell basepath
const WinCmdPath string = "cmd.exe"

// DefUserAgent - Default user agent string
const DefUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36"

// WebTimeout is the timeout for web requests
const WebTimeout = 5

// ShodanSearchTemplateURL is the Shodan URL Search template
const ShodanSearchTemplateURL = "https://www.shodan.io/search?query="

// GoogleSearchTemplateURL is the Google URL search template
const GoogleSearchTemplateURL = "https://www.google.com/search?q="

// GithubSearchTemplateURL is the Github URL search template
const GithubSearchTemplateURL = "https://github.com/search?q="

// execCheck is generally used to execute particular commands
func execCheck(target Target, checkID string, checkDetails CheckStruct,
	outfolder string, browserPath string, extensionsToExclude string,
	overwriteOutfiles bool) {
	
	//log.Printf("[v] checkDetails: %+v", checkDetails)
	checkType := checkDetails.Type

	log.Printf("[*] Executing checkID: %s of type: %s on target: %+v\n",
		checkID, checkType, target)
	if checkType == "cmd" {
		execCmd(target, checkID, checkDetails, outfolder, overwriteOutfiles)
	} else if checkType == "aws" || checkType == "awscli" {
		execAWSCLICmd(target, checkID, checkDetails, outfolder, overwriteOutfiles)
	} else if checkType == "gcloud" {
		execGCloudCmd(target, checkID, checkDetails, outfolder, overwriteOutfiles)
	} else if checkType == "bq" {
		execBQCmd(target, checkID, checkDetails, outfolder, overwriteOutfiles)
	} else if checkType == "webrequest" {
		execWebRequest(target, checkID, checkDetails, outfolder, overwriteOutfiles)
	} else if checkType == "grep" {
		execGrepSearch(target, checkID, checkDetails, extensionsToExclude,
			outfolder, overwriteOutfiles)
	} else if checkType == "find" {
		execFindSearch(target, checkID, checkDetails, outfolder, overwriteOutfiles)
	} else if checkType == "browser" || checkType == "webbrowser" || 
	  checkType == "webbrowse" || checkType == "browse" {
		execURLInBrowser(target, checkID, checkDetails, browserPath)
	} else if checkType == "shodan" {
		execShodanSearchInBrowser(target, checkID, checkDetails, browserPath)
	} else if checkType == "google" {
		execGoogleSearchInBrowser(target, checkID, checkDetails, browserPath)
	} else if checkType == "github" {
		execGithubSearchInBrowser(target, checkID, checkDetails, browserPath)
	} else if checkType == "notes" || checkType == "" {
		// Do nothing with notes
	} else {
		log.Printf("[-] Unknown check: %s of type: %s\n", checkID, checkType)
	}

}

// execCodeSearch is used to run the search in the folder via grep for specific
// keywords
func execGrepSearch(target Target, checkID string, checkDetails CheckStruct, 
	extensionsToExclude string, outfolder string, overwriteOutfiles bool) {

	keywords := checkDetails.Keywords
	outfile := checkDetails.Outfile
	writeToOutfileFlag := checkDetails.WriteToOutfile

	// Build the extensions command to exclude
	extensionsToExcludeCmd := ""
	for _, extn := range strings.Split(extensionsToExclude, ",") {
		extnWS := strings.TrimSpace(extn)
		extensionsToExcludeCmd += fmt.Sprintf(" --exclude=*.%s", extnWS)
	}

	cmdTemplate := "grep -A1 -B1 --color=always -rniE \"{keyword}\" \"{folder}\""
	cmdTemplate += extensionsToExcludeCmd

	for _, keyword := range keywords {
		cmdToExec := strings.ReplaceAll(cmdTemplate, "{keyword}", keyword)
		joinedCmds := subTargetParams(cmdToExec, target)
		totalOut := eCmd([]string{joinedCmds}, "")

		outfile = generateOutfile(checkID, writeToOutfileFlag,
			outfile, target)
		writeToOutfile(outfile, overwriteOutfiles, outfolder, totalOut, target)
	}
}

// execURLInBrowser opens URL(s) in a browser
func execURLInBrowser(target Target, checkID string, checkDetails CheckStruct,
	browserPath string) {

	urls := checkDetails.Urls

	for _, url := range urls {

		urlToOpen := subTargetParams(url, target)
		openURLInBrowser(urlToOpen, browserPath)
	}
}

// execShodanSearchInBrowser opens URL(s) in a browser
func execShodanSearchInBrowser(target Target, checkID string, checkDetails CheckStruct, 
	browserPath string) {

	searchesQueries := checkDetails.Searches

	for _, searchQuery := range searchesQueries {

		// Prepare a shodan search URL and open in browser
		searchURL := ShodanSearchTemplateURL + searchQuery
		urlToOpen := subTargetParams(searchURL, target)
		openURLInBrowser(urlToOpen, browserPath)
	}
}

// execGoogleSearchInBrowser opens URL(s) in a browser
func execGoogleSearchInBrowser(target Target, checkID string, checkDetails CheckStruct, 
	browserPath string) {

	searchesQueries := checkDetails.Searches

	for _, searchQuery := range searchesQueries {

		// Prepare a Google search URL and open in browser
		searchURL := GoogleSearchTemplateURL + searchQuery
		urlToOpen := subTargetParams(searchURL, target)
		openURLInBrowser(urlToOpen, browserPath)
	}
}

// execGithubSearchInBrowser opens URL(s) in a browser
func execGithubSearchInBrowser(target Target, checkID string, checkDetails CheckStruct, 
	browserPath string) {

	searchesQueries := checkDetails.Searches

	for _, searchQuery := range searchesQueries {

		// Prepare a Google search URL and open in browser
		searchURL := GithubSearchTemplateURL + searchQuery
		urlToOpen := subTargetParams(searchURL, target)
		openURLInBrowser(urlToOpen, browserPath)
	}
}

// execCodeSearch is used to run the search on folder for specific files
func execFindSearch(target Target, checkID string, checkDetails CheckStruct, 
	outfolder string, overwriteOutfiles bool) {

	files := checkDetails.Files
	outfile := checkDetails.Outfile
	writeToOutfileFlag := checkDetails.WriteToOutfile

	cmdTemplate := "find \"{folder}\" -ipath \"*{file}\""

	for _, file := range files {
		cmdToExec := strings.ReplaceAll(cmdTemplate, "{file}", file)
		joinedCmds := subTargetParams(cmdToExec, target)
		totalOut := eCmd([]string{joinedCmds}, "")

		outfile = generateOutfile(checkID, writeToOutfileFlag,
			outfile, target)
		writeToOutfile(outfile, overwriteOutfiles, outfolder, totalOut, target)
	}
}

// execCmd is used to execute shell commands and return the results
func execCmd(target Target, checkID string, checkDetails CheckStruct, 
	outfolder string, overwriteOutfiles bool) {

	// Read the necessary variables to execute
	cmdDir := checkDetails.CmdDir
	cmds := checkDetails.Cmds
	regex := checkDetails.Regex
	noregex := checkDetails.NoRegex
	alertOnMissing := checkDetails.AlertOnMissing
	outfile := checkDetails.Outfile
	writeToOutfileFlag := checkDetails.WriteToOutfile

	// Substitue target params in the command
	var subCmds []string
	for _, cmd := range cmds {
		subCmd := subTargetParams(cmd, target)
		subCmds = append(subCmds, subCmd)
	}
	// Execute the command to write to output
	totalOut := eCmd(subCmds, cmdDir)

	// If matching regex found, then print the result
	if shouldNotify(totalOut, regex, noregex, alertOnMissing) {
		generateStdOutNotification(checkDetails.Type, checkID, target.Target)
	} else {
		outfile = generateOutfile(checkID, writeToOutfileFlag,
			outfile, target)
		writeToOutfile(outfile, overwriteOutfiles, outfolder, totalOut, target)
	}

}

// execAWSCmd is used to execute AWSCLI commands and return the results
func execAWSCLICmd(target Target, checkID string, checkDetails CheckStruct, 
	outfolder string, overwriteOutfiles bool) {

	// Read the necessary variables to execute
	cmdDir := checkDetails.CmdDir
	cmds := checkDetails.Cmds
	regex := checkDetails.Regex
	noregex := checkDetails.NoRegex
	alertOnMissing := checkDetails.AlertOnMissing
	outfile := checkDetails.Outfile
	writeToOutfileFlag := checkDetails.WriteToOutfile

	// Convert commands to AWS Commands
	var awsCmds []string
	for _, cmd := range cmds {
		awsCmd := subTargetParams("aws "+cmd+" --profile={aws_profile} --region={aws_region}",
			target)
		awsCmds = append(awsCmds, awsCmd)
	}

	// Execute the command to write to output
	totalOut := eCmd(awsCmds, cmdDir)

	// If matching regex found, then print the result
	if shouldNotify(totalOut, regex, noregex, alertOnMissing) {
		generateStdOutNotification(checkDetails.Type, checkID, target.Target)
	} else {
		outfile = generateOutfile(checkID, writeToOutfileFlag, outfile, target)
		writeToOutfile(outfile, overwriteOutfiles, outfolder, totalOut, target)
	}

}

// execGCloudCmd is used to execute shell commands with gcloud and return results
func execGCloudCmd(target Target, checkID string, checkDetails CheckStruct, 
	outfolder string, overwriteOutfiles bool) {

	// Read the necessary variables to execute
	cmdDir := checkDetails.CmdDir
	cmds := checkDetails.Cmds
	regex := checkDetails.Regex
	noregex := checkDetails.NoRegex
	alertOnMissing := checkDetails.AlertOnMissing
	outfile := checkDetails.Outfile
	writeToOutfileFlag := checkDetails.WriteToOutfile
	disableYes := checkDetails.DisableYes

	// Set the default configuration before running the gcloud command
	gcloudConfigCmdsTemplate := []string{
		"config set project {gcp_project}",
		"config set compute/region {gcp_region}",
		"config set compute/zone {gcp_zone}",
		"config set account {gcp_account}",
	}
	var gcloudConfigCmds []string
	for _, cmd := range gcloudConfigCmdsTemplate {
		gcloudConfigCmd := subTargetParams("gcloud "+cmd, target)
		gcloudConfigCmds = append(gcloudConfigCmds,  gcloudConfigCmd)
	}
	eCmd(gcloudConfigCmds, cmdDir)

	// Convert actual commands to run to GCP Commands
	var gcloudCmds []string
	for _, cmd := range cmds {
		gcloudCmd := ""
		if disableYes {
			gcloudCmd = subTargetParams("gcloud "+cmd, target)
		} else {
			gcloudCmd = subTargetParams("yes | gcloud "+cmd, target)
		}
		gcloudCmds = append(gcloudCmds, gcloudCmd)
	}

	// Execute the command to write to output
	totalOut := eCmd(gcloudCmds, cmdDir)

	// If matching regex found, then print the result
	if shouldNotify(totalOut, regex, noregex, alertOnMissing) {
		generateStdOutNotification(checkDetails.Type, checkID, target.Target)
	} else {
		outfile = generateOutfile(checkID, writeToOutfileFlag, outfile, target)
		writeToOutfile(outfile, overwriteOutfiles, outfolder, totalOut, target)
	}

}

// execBQCmd is used to execute shell commands with gcloud and return results
func execBQCmd(target Target, checkID string, checkDetails CheckStruct, 
	outfolder string, overwriteOutfiles bool) {

	// Read the necessary variables to execute
	cmdDir := checkDetails.CmdDir
	cmds := checkDetails.Cmds
	regex := checkDetails.Regex
	noregex := checkDetails.NoRegex
	alertOnMissing := checkDetails.AlertOnMissing
	outfile := checkDetails.Outfile
	writeToOutfileFlag := checkDetails.WriteToOutfile

	// Convert commands to AWS Commands
	var bqCmds []string
	for _, cmd := range cmds {
		bqCmd := subTargetParams("bq "+cmd, target)
		bqCmds = append(bqCmds, bqCmd)
	}

	// Execute the command to write to output
	totalOut := eCmd(bqCmds, cmdDir)

	// If matching regex found, then print the result
	if shouldNotify(totalOut, regex, noregex, alertOnMissing) {
		generateStdOutNotification(checkDetails.Type, checkID, target.Target)
	} else {
		outfile = generateOutfile(checkID, writeToOutfileFlag, outfile, target)
		writeToOutfile(outfile, overwriteOutfiles, outfolder, totalOut, target)
	}

}

// execWebRequest is used to execute web requests on a specific target given the
// relevant check
func execWebRequest(target Target, checkID string, checkDetails CheckStruct, 
	outfolder string, overwriteOutfiles bool) {
	// Read vars for processing
	urls := checkDetails.Urls
	httpMethod := checkDetails.HTTPMethod
	regex := checkDetails.Regex
	noregex := checkDetails.NoRegex
	alertOnMissing := checkDetails.AlertOnMissing
	mheaders := checkDetails.Headers
	mbody := checkDetails.Body
	outfile := checkDetails.Outfile
	writeToOutfileFlag := checkDetails.WriteToOutfile

	totalOut := ""

	// Create the restyClient for making web requests in this thread
	restyClient := resty.New()

	// Set timeout and ignore verification of certificates
	restyClient.SetTimeout(time.Duration(WebTimeout) * time.Second)
	restyClient.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: true })

	for _, urlToCheck := range urls {

		// Determine if HTTP method is supported
		httpMethod := strings.ToUpper(httpMethod)
		if httpMethod == "" {
			httpMethod = "GET"
		}

		// Currently, we only support specific HTTP methods
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
			log.Println("[-] ", errResty)
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
			
			//fmt.Println(requestOut)

			// If matching regex found, then print the result
			if shouldNotify(requestOut, regex, noregex, alertOnMissing) {
				generateStdOutNotification(checkDetails.Type, checkID, urlToCheckSub)
			} else {
				outfile = generateOutfile(checkID, writeToOutfileFlag, outfile, 
					target)
				writeToOutfile(outfile, overwriteOutfiles, outfolder, requestOut, target)
			}

			// Append to full output to be used later (if necessary)
			totalOut += requestOut
		}
	}
}

// execNotes is used to print the notes given the target and the method details
/*func execNotes(target *Target, method *MethodStruct) {
	// Read the necessary variables to print notes
	notes := method.Notes
	log.Println("[!] Notes:")
	notesToPrint := strings.Split(notes, "\n")
	for _, note := range notesToPrint {
		noteToPrint := subTargetParams(note, *target)
		log.Println("[!] " + noteToPrint)
	}
}*/