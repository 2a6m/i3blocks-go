package main

import(
	"flag"
	"fmt"
	"os"
	"syscall"
)

type DiskStatus struct {
	All float64
	Used float64
	Free float64
}

const (
	B = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

func diskUsage(path string) (disk DiskStatus, err error) {
	fs := syscall.Statfs_t{}
	err = syscall.Statfs(path, &fs)
	if err != nil {
		return disk, err
	}
	disk.All = float64(fs.Blocks * uint64(fs.Bsize))/float64(GB)
	disk.Free = float64(fs.Bfree * uint64(fs.Bsize))/float64(GB)
	disk.Used = disk.All - disk.Free
	return disk, nil
}

func main() {

	// Set display texts to defaults.
	var output string
	var fullText string = "error"
	var shortText string = "error"
	var color string = "#ff0000"

	// Set flags
	pathCmd := flag.String("path", "/", "set the path of the mountpoint for the disk")
	flag.Parse()

	// Get disk usage from the given path
	disk, err := diskUsage(*pathCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[i3blocks-go] Failed to read battery files: %s", err.Error())
		fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n", fullText, shortText, color)
		os.Exit(0)
	}

	output = fmt.Sprintf("%.2f / %.2f GB", disk.Used, disk.All)
	fullText = output
	shortText = output
	if (disk.Used / disk.All) > 0.8 {
		color = "#ff0000"
	} else {
		color = "#ffffff"
	}

	fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n", fullText, shortText, color)
	os.Exit(0)
}

