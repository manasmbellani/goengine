package main

import (
	"fmt"
	"log"
	"os"
)

// BrowserPaths is a list of possible browser paths
var BrowserPaths = []string{
"/Applications/Brave Browser.app/Contents/MacOS/Brave Browser",
"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
"/Applications/Firefox.app/Contents/MacOS/firefox",
"/Applications/Safari.app/Contents/MacOS/Safari",
}

// locateBrowserPath is used to locate the appropriate browser path
func locateBrowserPath() string{
	browserPath := ""
	for _, bp := range BrowserPaths {
		_, err := os.Stat(bp)
		if err == nil {
			browserPath = bp
			break
		} 
	}
	return browserPath
}

// openURLInBrowser opens a URL in the specified browser
func openURLInBrowser (url string, browserPath string, target Target) {
	// Cannot open URL if browser not found
	if browserPath == "" {
		log.Println("[-] Browser not found to open URL: ", url)
	} else {
		cmdToExec := fmt.Sprintf("\"%s\" \"%s\" 2>/dev/null 1>/dev/null &", 
			browserPath, url)
		eCmd([]string{cmdToExec}, "", target)
	}

}