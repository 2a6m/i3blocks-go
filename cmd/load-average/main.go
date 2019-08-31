package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"

	"io/ioutil"
)

const sysfs = "/proc/loadavg"

func load1() (float64, error){
	loadRaw, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return -1, err
	}
	// Remove surrounding space and split at inner spaces.
	loadStrings := strings.Split(strings.TrimSpace(string(loadRaw)), " ")
	loadFloat, err := strconv.ParseFloat(loadStrings[0], 64)
	if err != nil {
		return -1, err
	}
	return loadFloat, err
}

func load5() (float64, error){
	loadRaw, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return -1, err
	}
	// Remove surrounding space and split at inner spaces.
	loadStrings := strings.Split(strings.TrimSpace(string(loadRaw)), " ")
	loadFloat, err := strconv.ParseFloat(loadStrings[1], 64)
	if err != nil {
		return -1, err
	}
	return loadFloat, err
}

func load15() (float64, error){
	loadRaw, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return -1, err
	}
	// Remove surrounding space and split at inner spaces.
	loadStrings := strings.Split(strings.TrimSpace(string(loadRaw)), " ")
	loadFloat, err := strconv.ParseFloat(loadStrings[2], 64)
	if err != nil {
		return -1, err
	}
	return loadFloat, err
}
func main() {

	// Set display texts to defaults.a
	var output string
	var fullText string = "error"
	var shortText string = "error"
	var color string = "#ff0000"

	// Read current load average information from kernel
	// pseudo-file-system mounted at /proc.
	load, err := load5()
	if err != nil {

		// Write an error to STDERR, fallback display values
		// to STDOUT and exit with failure code.
		fmt.Fprintf(os.Stderr, "[i3blocks-go] Failed to read load average file: %s", err.Error())
		fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n", fullText, shortText, color)
		os.Exit(0)
	}


	// Depending on length of display text, construct
	// final output string.
	output = fmt.Sprintf("CPU: %.2f", load)
	fullText = output
	shortText = output
	if load > 0.8 {
		color = "#ff0000"
	} else {
		color = "#ffffff"
	}

	// Write out gathered information to STDOUT.
	fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n", fullText, shortText, color)
	os.Exit(0)
}
