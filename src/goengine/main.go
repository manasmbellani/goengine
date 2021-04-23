// Golang script to print greetings as an example for simple goroutines pattern
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"

	glob "github.com/ganbarodigital/go_glob"
	"github.com/go-resty/resty/v2"
)

// DefPort is the default port
const DefPort string = "443"

// DefProtocol is the default protocol to use
const DefProtocol string = "https"

// Extensions to exclude
const DefExtensionsToExclude string = "jpg,svg,png,bmp,css"

func printGreetingsWorkers(names *chan string, greeting string, numThreads int,
	wg *sync.WaitGroup) {
	// Execute workers to print greetings

	// Start threads to complete task
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for name := range *names {
				// Print greeting
				fmt.Printf("%s %s\n", greeting, name)
			}
		}()
	}
}

// normalizeTargetWorkers starts goroutines to convert raw targets to individual
// target parts
func normalizeTargetWorkers(targets *[]Target, rawTargets chan string,
	outfolder string, numThreads int, wg *sync.WaitGroup) {
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for rawTarget := range rawTargets {
				// Normalize the target
				var target Target
				normalizeTarget(rawTarget, &target, outfolder)

				// Print the information about the target
				log.Printf("Added target: %+v to targets", target)
				*targets = append(*targets, target)

			}
		}()
	}
}

// parseCheckFiles is used to parse the check file one by one
func parseCheckFiles(checkFiles []string, allChecks map[string]CheckStruct) {
	for _, checkFile := range checkFiles {
		checkStruct := parseCheckFile(checkFile)
		allChecks[checkFile] = checkStruct
	}
}

// Parse the Checks file to a structure that we can read from
func parseCheckFile(checkFile string) CheckStruct {
	var checksFileStruct CheckStruct
	yamlFile, err := ioutil.ReadFile(checkFile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &checksFileStruct)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return checksFileStruct
}

// execChecksWorkers executes the checks
func execChecksWorkers(checksToExec chan CheckToExec, restyClient *resty.Client,
	numThreads int, outfolder string, browserPath string, extensionsToExclude string,
	wg *sync.WaitGroup) {
	for i := 0; i < numThreads; i++ {
		log.Printf("[*] Launching worker: %d for execChecksWorker\n", i)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for checkToExec := range checksToExec {

				// Execute the method based on the target
				target := checkToExec.Target
				checkDetails := checkToExec.CheckDetails
				checkID := checkToExec.CheckID
				execCheck(target, checkID, checkDetails, outfolder, browserPath,
					extensionsToExclude)
			}
		}()
	}
}

// prepareChecksToExecWorkers is used to prepare the checks to execute and also
// determine whether a check should be executed or not (based on glob user input)
func prepareChecksToExecWorkers(allChecks map[string]CheckStruct,
	targets []Target, checksToExec chan CheckToExec) {

	// Loop through each check and determine if check needs to be exec
	for checkID, checkDetails := range allChecks {
		for _, t := range targets {
			var checkToExec CheckToExec
			checkToExec.CheckID = checkID
			checkToExec.CheckDetails = checkDetails
			checkToExec.Target = t			
			log.Printf("[*] Added check ID: %s for target: %s to checksToExec\n",
				checkID, t)
			checksToExec <- checkToExec
		}
	}
}

// shouldExecCheck determines if the check should be executed
func shouldExecCheck(inp string, match string) bool {
	if match == "all" {
		return true
	}
	return containsGlobPattern(inp, match)
}

// containsGlobPattern checks if the input contains matching string
func containsGlobPattern(inp string, match string) bool {

	// Add * for 'contains' matching via glob
	if !strings.Contains(match, "*") {
		match = fmt.Sprintf("*%s*", match)
	}

	// Check if match found in input
	myGlob := glob.NewGlob(match)
	success, _ := myGlob.Match(inp)
	return success
}

// getCheckFiles gets the list of checks files
func getCheckFiles(checksFolder string) []string {
	var checkFiles []string
	fi, err := os.Stat(checksFolder)
	if os.IsNotExist(err) {
		log.Printf("[-] Checks folder: %s does not exist\n", checksFolder)
	} else {
		mode := fi.Mode()
		if mode.IsRegular() {
			// Folder is actually a single file, let's parse that
			checkFile := checksFolder
			checkFiles = append(checkFiles, checkFile)
		} else if mode.IsDir() {
			// Parse all the files in the folder
			err = filepath.Walk(checksFolder,
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					fi, err := os.Stat(path)
					if err != nil {
						log.Fatalf("Error getting path info: %s\n", err.Error())
						return err
					}

					// Add the files for monitoring 
					mode := fi.Mode()
					if mode.IsRegular() {
						checkFiles = append(checkFiles, path)
					}
					return nil
				},
			)
			if err != nil {
				log.Fatalf("Error walking the directory: %s. Err: %s\n", checksFolder, 
					err.Error())
			}
		} else {
			log.Fatalf("Path: %s is neither file, nor directory", checksFolder)
		}
	}

	return checkFiles
}

// getCheckFilesToExec determines whether a check should be executed
func getCheckFilesToExec(allCheckFiles []string, checkIDsToExec string) []string {
	var checkFilesToExec []string
	for _, checkFile := range allCheckFiles {
		log.Printf("[*] Appended checkfile: %s for processing", checkFile)
		if shouldExecCheck(checkFile, checkIDsToExec) {
			checkFilesToExec = append(checkFilesToExec, checkFile)
		}
	}
	return checkFilesToExec
}

func main() {
	var checksFolder string
	var numThreads int
	var numThreadsNT int
	var checkIDsToExec string
	var outfolder string
	var extensionsToExclude string
	var quiet bool

	flag.StringVar(&checksFolder, "f", "vulnreview.yaml", "Checks File in YAML")
	flag.StringVar(&checkIDsToExec, "c", "all", "Checks to execute")
	flag.IntVar(&numThreads, "numThreads", 50,
		"Number of threads for vuln scanning")
	flag.IntVar(&numThreadsNT, "numThreadsNT", 4,
		"Number of threads for normalization of targets")
	flag.StringVar(&outfolder, "outfolder", "/opt/dockershare/goengine",
		"Folder where the outfiles are written")
	flag.StringVar(&extensionsToExclude, "ee", DefExtensionsToExclude,
		"Extensions to exclude when performing grep searches")
	flag.BoolVar(&quiet, "q", false,
		"Execute in quiet mode so no verbose messages are printed")
	flag.Parse()

	// Signature file should be found
	if _, err := os.Stat(checksFolder); os.IsNotExist(err) {
		log.Fatalf("Checks Files/folder: %s does not exist\n", checksFolder)
	}

	// Quiet mode
	if quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	// Determine the local browser
	browserPath := locateBrowserPath()
	if browserPath == "" {
		log.Printf("[-] Browser not found to run browser checks")
	}

	var allChecks map[string]CheckStruct

	// Parse the checks file
	allChecksFiles := getCheckFiles(checksFolder)

	// Determine Checks File to execute
	checksFiles := getCheckFilesToExec(allChecksFiles, checkIDsToExec)
	
	// Parse the checks file
	parseCheckFiles(checksFiles, allChecks)

	// Create sync group for normalization of the targets, perparing checks to
	// execute and executing the check
	var wgNT sync.WaitGroup
	var wgEC sync.WaitGroup

	// targets to parse
	rawTargets := make(chan string)

	// Parse targets into this structs list
	var targets []Target

	// Create a resty client and initialise it
	restyClient := resty.New()

	// Keep a list of checks to execute
	checksToExec := make(chan CheckToExec)

	// Parse the targets
	normalizeTargetWorkers(&targets, rawTargets, outfolder, numThreadsNT, &wgNT)

	// Read assets to process from STDIN input
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line != "" {
			rawTargets <- line
		}
	}

	// No more names, close channel and wait for print greeting workers to end
	close(rawTargets)
	wgNT.Wait()

	// Start workers to execute the checks
	execChecksWorkers(checksToExec, restyClient, numThreads, outfolder,
		browserPath, extensionsToExclude, &wgEC)

	// Prepare a list of the relevant checks to execute for each target
	prepareChecksToExecWorkers(allChecks, targets, checksToExec)

	close(checksToExec)

	wgEC.Wait()

}
