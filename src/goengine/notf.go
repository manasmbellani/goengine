package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// NewLineReplacement is char to replace new lines with for regex search
const NewLineReplacement = "|"

// OutfilePrefix is the prefix for output files when generating them
const OutfilePrefix = "out"

// OutfileExtn is the extension for the outfile
const OutfileExtn = "txt"

// Get filename without extension
func fileNameWithoutExtension(filePath string) string {
	fileName := filepath.Base(filePath)
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

// shouldNotify is used to send notification based on input line and regex
func shouldNotify(out string, regex string, noregex string, 
	alertOnMissing bool) bool {
	found := false
	if regex != "" {
		outWithoutNewLines := strings.ReplaceAll(out, "\n", NewLineReplacement)
		outWithoutNewLines = strings.ReplaceAll(outWithoutNewLines, "\r", NewLineReplacement)
		foundRegexMatch, _ := regexp.MatchString(regex, outWithoutNewLines)
		var foundNoRegexMatch bool
		if noregex != "" {
			foundNoRegexMatch, _ = regexp.MatchString(noregex, outWithoutNewLines)
		}
		if alertOnMissing {
			found = !(foundRegexMatch && !foundNoRegexMatch)
		} else {
			found = (foundRegexMatch && !foundNoRegexMatch)
		}
	}
	return found
}

// generateStdOutNotification writes a message on STDOUT
func generateStdOutNotification(checkType, checkID, asset string) {
	// checkID is currently a full file path so we replace it with just the 
	// filename 
	checkIDShort := fileNameWithoutExtension(checkID)
	fmt.Printf("[%s] %s\n", checkIDShort, asset)
}

func writeToOutfile(outfile string, overwriteOutfiles bool, outfolder string, 
	out string, target Target) {
	// Append or overwrite the output in the outfolder/outfile
	if outfile != "" {

		outfileSub := subTargetParams(outfile, target)
		outfileFullPath := filepath.Join(outfolder, outfileSub)
		log.Printf("[*] Writing results to outfile: %s\n", outfileFullPath)
		mode := os.O_APPEND|os.O_CREATE|os.O_WRONLY
		if overwriteOutfiles {
			mode = os.O_RDWR|os.O_CREATE|os.O_TRUNC
		}
		f, err := os.OpenFile(outfileFullPath, mode, 0644)
		if err != nil {
			log.Println("[*] Error opening file: ", outfileFullPath, err)
		}
		defer f.Close()
		if _, err := f.WriteString(out); err != nil {
			log.Println("[-] " + err.Error())
		}
	}
}

// generateOutfile is used to generate the output file name
func generateOutfile(checkID string, writeToOutfile bool,
	outfileCheck string, target Target) string {

	outfile := ""
	
	// Get the name of the check for preparing filename
	checkName := fileNameWithoutExtension(checkID)

	if outfileCheck != "" {
		outfile = outfileCheck
	} else if writeToOutfile {
		protocol := target.Protocol
		if protocol == "folder" {
			// Build the file name replacing disallowed characters with '_'
			folder_name := strings.ReplaceAll(target.Folder, "/", "_")
			folder_name = strings.ReplaceAll(folder_name, "\\", "_")
			outfile = fmt.Sprintf("%s-%s-%s.%s", OutfilePrefix, checkName,
				folder_name, OutfileExtn)
		} else if protocol == "aws" {
			outfile = fmt.Sprintf("%s-%s-%s-%s.%s", OutfilePrefix, checkName,
				target.AWSProfile, target.AWSRegion, OutfileExtn)
		} else if protocol == "gcp" || protocol == "gcloud" {
			outfile = fmt.Sprintf("%s-%s-%s-%s-%s-%s.%s", OutfilePrefix, checkName,
				target.GCPAccount, target.GCPProject, target.GCPRegion, target.GCPZone,
				OutfileExtn)
		} else {
			outfile = fmt.Sprintf("%s-%s-%s.%s", OutfilePrefix, checkName,
				target.Host, OutfileExtn)
		}
	} else {
		outfile = ""
	}

	return outfile
}
