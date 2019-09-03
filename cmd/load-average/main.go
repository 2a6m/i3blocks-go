package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"strconv"
	"runtime"

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

func main() {

	// Set display texts to defaults.a
	var output string
	var fullText string = "error"
	var shortText string = "error"
	var color string = "#ff0000"

	// flags
	var timeFlag = flag.Int("time", 1, "Set -time to set the cpu load average wanted [0 - for 1 min, 1 - for 5 min, 2 - for 15 min]")
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

	// Depending on length of display text, construct
	// final output string.
	l = l / float64(runtime.NumCPU())
	output = fmt.Sprintf("%.2f", l)
	fullText = output
	shortText = output
	if l >= 1 {
		color = "#ff0000"
	} else if l >= 0.75 {
		color = "#ff8800"
	} else if l >= 0.5 {
		color = "#ffff00"
	} else if l >= 0.25 {
		color = "#88ff00"
	} else {
		color = "#00ff00"
	}

	// Write out gathered information to STDOUT.
	fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n", fullText, shortText, color)
	os.Exit(0)
}
