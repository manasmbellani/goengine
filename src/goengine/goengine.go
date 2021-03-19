// Golang script to print greetings as an example for simple goroutines pattern
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"

	glob "github.com/ganbarodigital/go_glob"
	"github.com/go-resty/resty"
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
	numThreads int, wg *sync.WaitGroup) {
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for rawTarget := range rawTargets {
				// Normalize the target
				var target Target
				normalizeTarget(rawTarget, &target)

				// Print the information about the target
				log.Printf("Added target: %+v to targets", target)
				*targets = append(*targets, target)

			}
		}()
	}
}

// Parse the Checks file to a structure that we can read from
func parseChecksFile(checksFile string, allchecks *[]CheckStruct) {
	var checksFileStruct ChecksFileStruct
	yamlFile, err := ioutil.ReadFile(checksFile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &checksFileStruct)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	*allchecks = checksFileStruct.Checks
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
				method := checkToExec.Method
				checkID := checkToExec.CheckID
				methodID := checkToExec.MethodID
				execMethod(target, checkID, methodID, method, outfolder, 
					browserPath, extensionsToExclude)
			}
		}()
	}
}

// prepareChecksToExecWorkers is used to prepare the checks to execute and also
// determine whether a check should be executed or not (based on glob user input)
func prepareChecksToExecWorkers(allChecks []CheckStruct,
	targets []Target, checksToExec chan CheckToExec,
	checkIDsToExec string, methodIDsToExec string) {

	// Loop through each check and determine if check needs to be exec
	for _, check := range allChecks {
		checkID := check.ID
		if shouldExecCheck(checkID, checkIDsToExec) {

			// Check if the method needs to be executed
			methods := check.Methods
			for _, method := range methods {
				methodID := method.ID
				if shouldExecCheck(methodID, methodIDsToExec) {
					// Add target and method info as a check to execute
					// listing
					for _, t := range targets {
						var checkToExec CheckToExec
						checkToExec.CheckID = checkID
						checkToExec.MethodID = methodID
						checkToExec.Target = t
						checkToExec.Method = method
						log.Printf("[*] Added check: %s, method: %s for target: %s to checksToExec\n",
							checkID, methodID, t)
						checksToExec <- checkToExec
					}
				}
			}
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

func main() {
	var checksFile string
	var numThreads int
	var numThreadsNT int
	var checkIDsToExec string
	var methodIDsToExec string
	var outfolder string
	var extensionsToExclude string
	var quiet bool
	flag.StringVar(&checksFile, "f", "vulnreview.yaml", "Checks File in YAML")
	flag.StringVar(&checkIDsToExec, "c", "all", "Checks to execute")
	flag.StringVar(&methodIDsToExec, "m", "all", "Methods to execute")
	flag.IntVar(&numThreads, "numThreads", 10,
		"Number of threads for vuln scanning")
	flag.IntVar(&numThreadsNT, "numThreadsNT", 2,
		"Number of threads for normalization of targets")
	flag.StringVar(&outfolder, "outfolder", "/opt/dockershare/goengine",
		"Folder where the outfiles are written")
	flag.StringVar(&extensionsToExclude, "ee", DefExtensionsToExclude, 
		"Extensions to exclude when performing grep searches")
	flag.BoolVar(&quiet, "q", false,
		"Execute in quiet mode so no verbose messages are printed")
	flag.Parse()

	// Signature file should be found
	if _, err := os.Stat(checksFile); os.IsNotExist(err) {
		log.Fatalf("Checks File: %s does not exist\n", checksFile)
	}

	// Quiet mode
	if quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	// Parse the checks file and end
	var allchecks []CheckStruct
	parseChecksFile(checksFile, &allchecks)

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
	normalizeTargetWorkers(&targets, rawTargets, numThreadsNT, &wgNT)

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

	// Determine the local browser
	browserPath := locateBrowserPath()
	if browserPath == "" {
		log.Printf("[-] Browser not found to run browser checks")
	}

	// Start workers to execute the checks
	execChecksWorkers(checksToExec, restyClient, numThreads, outfolder, 
		browserPath, extensionsToExclude, &wgEC)

	// Prepare a list of the relevant checks to execute for each target
	prepareChecksToExecWorkers(allchecks, targets, checksToExec, checkIDsToExec,
		methodIDsToExec)

	close(checksToExec)

	wgEC.Wait()

}
