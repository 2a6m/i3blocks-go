package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"io/ioutil"
	"path/filepath"
	"os/exec"
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

func notify(title, msg string) (string, error) {
	cmd := exec.Command("notify-send", "--urgency=critical", title, msg)
	bout, err := cmd.Output()
	out := string(bout)
	if err != nil {
		return out, err
	}
	if out != "" {
		return out, nil
	}
	return  "", nil
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

	// flags
	var notificationFlag = flag.Bool("notification", false, "Pass -notification to send notifiaction when battery is low")
	var levelFlag = flag.Int("level", 15, "Set -level to set the notification trigger level")
	flag.Parse()

	// Read charging status information from kernel
	// pseudo-file-system mounted at /sys.
	b, err := battery()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[i3blocks-go] Failed to read battery files: %s", err.Error())
		fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n", fullText, shortText, color)
		os.Exit(0)
	}

	output = fmt.Sprintf("Bat: %.2f%%", b)
	fullText = output
	shortText = output
	if b < float64(*levelFlag) {
		if *notificationFlag {
			notify("Battery low", fullText)
		}
	} else {
		color = "#ffffff"
	}

	fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n", fullText, shortText, color)
	os.Exit(0)
}
