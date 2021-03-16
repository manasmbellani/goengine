package main

import (
	"os"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

// NewLineReplacement is char to replace new lines with for regex search
const NewLineReplacement string = "|"

// shouldNotify is used to send notification based on input line and regex
func shouldNotify(out string, regex string) bool {
	found := false
	if regex != "" {
		outWithoutNewLines := strings.ReplaceAll(out, "\n", NewLineReplacement)
		outWithoutNewLines = strings.ReplaceAll(outWithoutNewLines, "\r", NewLineReplacement)
		found, _ = regexp.MatchString(regex, outWithoutNewLines)
	}
	return found
}

func writeToOutfile(outfile string, outfolder string, out string,
	target Target) {

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
