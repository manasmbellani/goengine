package main

import (
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
