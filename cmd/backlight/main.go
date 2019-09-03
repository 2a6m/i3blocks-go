package main

import(
	"fmt"
	"os"
	"strconv"

	"io/ioutil"
	"path/filepath"
)

const sysfs = "/sys/class/backlight"

func readFloat(path, filename string) (float64, error) {
	str, err := ioutil.ReadFile(filepath.Join(path, filename))
	if err != nil {
		return 0, err
	}
	num, err := strconv.ParseFloat(string(str[:len(str)-1]), 64)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func getBacklightFiles() ([]string, error) {
	files, err := ioutil.ReadDir(sysfs)
	if err != nil {
		return nil, err
	}

	var bFiles []string
	for _, file := range files {
		path := filepath.Join(sysfs, file.Name())
		bFiles = append(bFiles, path)
	}
	return bFiles, err
}

func backlight() (float64, error){
	backlightFiles, err := getBacklightFiles()
	if err != nil {
		return 0, err
	}
	// there is only one file in backlight
	brightness, err := readFloat(backlightFiles[0], "actual_brightness")
	total, err := readFloat(backlightFiles[0], "max_brightness")

	if err != nil {
		return 0, err
	}
	p := (brightness / total)*100
	return p, nil
}

func main() {

	// Set display tests to defaults
	var output string
	var fullText string = "error"
	var shortText string = "error"
	var color string = "#ff0000"

	// Read brightness dtatud information from kernel
	// pseudo-file-system mounted at /sys
	b, err := backlight()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[i3block-go] Failed to read backlight files: %s", err.Error())
		fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n", fullText, shortText, color)
		os.Exit(0)
	}

	output = fmt.Sprintf("%.2f%%", b)
	fullText = output
	shortText = output
	color = "" // use default color define by i3blocks

	fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n", fullText, shortText, color)
	os.Exit(0)
}

