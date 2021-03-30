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

// shouldNotify is used to send notification based on input line and regex
func shouldNotify(out string, regex string, alertOnMissing bool) bool {
	found := false
	if regex != "" {
		outWithoutNewLines := strings.ReplaceAll(out, "\n", NewLineReplacement)
		outWithoutNewLines = strings.ReplaceAll(outWithoutNewLines, "\r", NewLineReplacement)
		foundMatch, _ := regexp.MatchString(regex, outWithoutNewLines)
		if alertOnMissing {
			found = !foundMatch
		} else {
			found = foundMatch
		}
	}
	return found
}

func writeToOutfile(outfile string, outfolder string, out string, target Target) {
	// Append the output to outfolder/outfile
	if outfile != "" {

		outfileSub := subTargetParams(outfile, target)
		outfileFullPath := filepath.Join(outfolder, outfileSub)
		log.Printf("[*] Writing results to outfile: %s\n", outfileFullPath)
		f, err := os.OpenFile(outfileFullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Error opening file: ", outfileFullPath, err)
		}
		defer f.Close()
		if _, err := f.WriteString(out); err != nil {
			log.Println(err)
		}
	}
}

// generateOutfile is used to generate the output file name
func generateOutfile(checkID string, methodID string, writeToOutfile bool,
	outfileMethod string, target Target) string {

	outfile := ""

	if outfileMethod != "" {
		outfile = outfileMethod
	} else if writeToOutfile {
		protocol := target.Protocol
		if protocol == "folder" {
			// Build the file name replacing disallowed characters with '_'
			folder_name := strings.ReplaceAll(target.Folder, "/", "_")
			folder_name = strings.ReplaceAll(folder_name, "\\", "_")
			outfile = fmt.Sprintf("%s-%s-%s-%s.%s", OutfilePrefix, checkID, methodID,
				folder_name, OutfileExtn)
		} else if protocol == "aws" {
			outfile = fmt.Sprintf("%s-%s-%s-%s-%s.%s", OutfilePrefix, checkID, methodID,
				target.AWSProfile, target.AWSRegion, OutfileExtn)
		} else if protocol == "gcp" || protocol == "gcloud" {
			outfile = fmt.Sprintf("%s-%s-%s-%s-%s-%s-%s.%s", OutfilePrefix, checkID, methodID,
				target.GCPAccount, target.GCPProject, target.GCPRegion, target.GCPZone,
				OutfileExtn)
		} else {
			outfile = fmt.Sprintf("%s-%s-%s-%s.%s", OutfilePrefix, checkID, methodID,
				target.Host, OutfileExtn)
		}
	} else {
		outfile = ""
	}

	return outfile
}
