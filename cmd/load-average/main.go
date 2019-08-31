package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"strconv"

	"io/ioutil"
)

const sysfs = "/proc/loadavg"

func load(i int) (float64, error){
	loadRaw, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return -1, err
	}
	// Remove surrounding space and split at inner spaces.
	loadStrings := strings.Split(strings.TrimSpace(string(loadRaw)), " ")
	loadFloat, err := strconv.ParseFloat(loadStrings[i], 64)
	if err != nil {
		return -1, err
	}
	return loadFloat, err
}

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

	// flags
	var timeFlag = flag.Int("time", 1, "Set -time to set the cpu load average wanted [0 - for 1 min, 1 - for 5 min, 2 - for 15 min]")
	var processFlag = flag.Bool("process", false, "Pass -process to show the number of runninn process")
	flag.Parse()

	// Read current load average information from kernel
	// pseudo-file-system mounted at /proc.
	var l float64
	var err error
	if *timeFlag <= 2 {
		l, err = load(*timeFlag)
	} else {
		l, err = load(1)
	}
	if err != nil {

		// Write an error to STDERR, fallback display values
		// to STDOUT and exit with failure code.
		fmt.Fprintf(os.Stderr, "[i3blocks-go] Failed to read load average file: %s", err.Error())
		fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n", fullText, shortText, color)
		os.Exit(0)
	}

	output = fmt.Sprintf("CPU: %.2f", l)
	var p int
	if *processFlag {
		p, err = process()
		output = fmt.Sprintf("CPU: %.2f  Process: %d", l, p)
	}

	// Depending on length of display text, construct
	// final output string.
	fullText = output
	shortText = output
	if l > 0.8 {
		color = "#ff0000"
	} else {
		color = "#ffffff"
	}

	// Write out gathered information to STDOUT.
	fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n", fullText, shortText, color)
	os.Exit(0)
}
