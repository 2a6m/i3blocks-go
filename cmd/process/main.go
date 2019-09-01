package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"

	"io/ioutil"
)

const sysfs = "/proc/loadavg"

func process() (int, error){
	loadRaw, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return -1, err
	}
	// Remove surrounding space and split at inner spaces.
	loadStrings := strings.Split(strings.TrimSpace(string(loadRaw)), " ")
	loadStrings = strings.Split(strings.TrimSpace(string(loadStrings[3])), "/")
	loadInt, err := strconv.Atoi(loadStrings[1])
	if err != nil {
		return -1, err
	}
	return loadInt, err
}

func main() {

	// Set display texts to defaults.a
	var output string
	var fullText string = "error"
	var shortText string = "error"
	var color string = "#ff0000"

	// Read current load average information from kernel
	// pseudo-file-system mounted at /proc.
	p, err := process()
	if err != nil {

		// Write an error to STDERR, fallback display values
		// to STDOUT and exit with failure code.
		fmt.Fprintf(os.Stderr, "[i3blocks-go] Failed to read load average file: %s", err.Error())
		fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n", fullText, shortText, color)
		os.Exit(0)
	}

	// Depending on length of display text, construct
	// final output string.
	output = fmt.Sprintf("Process: %d", p)
	fullText = output
	shortText = output
	color = "#ffffff"

	// Write out gathered information to STDOUT.
	fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n", fullText, shortText, color)
	os.Exit(0)
}
