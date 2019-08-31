package main

import (
	"fmt"
	"os"
	"strconv"

	"io/ioutil"
	"path/filepath"
)

const sysfs = "/sys/class/power_supply"

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

func isBattery(path string) bool {
	t, err := ioutil.ReadFile(filepath.Join(path,"type"))
	return err == nil && string(t) == "Battery\n"
}

func getBatteryFiles() ([]string, error) {
	files, err := ioutil.ReadDir(sysfs)
	if err != nil {
		return nil, err
	}

	var bFiles []string
	for _, file := range files {
		//fmt.Println(file.Name())
		path := filepath.Join(sysfs, file.Name())
		if isBattery(path) {
			bFiles = append(bFiles, path)
		}
	}
	return bFiles, err
}

func battery() (float64, error) {
	batteryFiles, err := getBatteryFiles()
	if err != nil {
		return 0, err
	}
	var charge float64
	var total float64
	for _, file := range batteryFiles {
		c, err := readFloat(file, "energy_now")
		t, err := readFloat(file, "energy_full_design")
		if err != nil {
			return 0, err
		}
		charge += c
		total += t
	}
	capacity := (charge / total)*100
	return capacity, nil
}


func main() {

	// Set display texts to defaults.
	var output string
	var fullText string = "error"
	var shortText string = "error"
	var color string = "#ff0000"

	// Read charging status information from kernel
	// pseudo-file-system mounted at /sys.
	b, err := battery()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[i3blocks-go] Failed to read battery files: %s", err.Error())
		fmt.Fprintf(os.Stdout, "%s\n%s\n%s", fullText, shortText, color)
		os.Exit(0)
	}

	output = fmt.Sprintf("bat: %.2f%%", b)
	if b < 20 {
		// if notification flag
	} else {
		color = "#000000"
	}
	fullText = output
	shortText = output

	fmt.Fprintf(os.Stdout, "%s\n%s\n%s", fullText, shortText, color)
	os.Exit(0)
}
