package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/mholt/archiver"
)

// Version of updater can be found by calling the executable on
// its own with just "--version"
var version = "1.0"

// Main usage:
//  Program should be called with the following arguments:
//  - Update download URL
//      - A tar.gz file of a single location
//  - Update folder target to overwrite, trailing slash required

func main() {
	// Checks for version argument
	if len(os.Args) == 2 {
		if os.Args[1] == "--version" {
			fmt.Println("Go-Updater Version " + version)
			os.Exit(1)
		}
	}

	// Check arguments
	if len(os.Args) < 3 {
		fmt.Println("Invalid arguments provided.")
		os.Exit(0)
	}

	// URL Checks
	var downloadURL string // Declared here to allow building with "go build"
	_, err := url.ParseRequestURI(os.Args[1])
	if err != nil {
		fmt.Println("Invalid URL argument provided.")
		fmt.Println(err.Error())
		os.Exit(0)
	}

	downloadURL = os.Args[1]

	// Checks update target location
	updateTarget, err := os.Stat(os.Args[2])
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Update target dir doesn't exist.")
		}

		fmt.Println("Error accessing update target dir")
		fmt.Println(err.Error())
		os.Exit(0)
	}

	if !updateTarget.IsDir() {
		fmt.Println("The update target dir appears to be a file.")
		os.Exit(0)
	}

	updateDir := os.Args[2]

	// Downloads the .gz to a temporary file
	err = os.MkdirAll("temp/", 0755)

	if err != nil {
		fmt.Println("Unable to create temporary directory.")
		fmt.Println(err.Error())
	}

	gzFile, err := os.Create("temp/update.tgz")

	if err != nil {
		fmt.Println("Can not create the temporary download file.")
		fmt.Println(err.Error())
		os.Exit(0)
	}

	resp, err := http.Get(downloadURL)

	if err != nil {
		fmt.Println("Can not download the update.")
		fmt.Println(err.Error())
		os.Exit(0)
	}

	defer resp.Body.Close()

	n, err := io.Copy(gzFile, resp.Body)
	if err != nil {
		fmt.Println("Can not download the update.")
		fmt.Println(err.Error())
		os.Exit(0)
	}

	fmt.Println("Successfully downloaded the update: " + strconv.Itoa(int(n)) + "bytes downloaded.")
	fmt.Println("Extracting update archive.")

	err = archiver.TarGz.Open("temp/update.tgz", "temp/update/")

	if err != nil {
		fmt.Println("Error extracting archive.")
		fmt.Println(err.Error())
		os.Exit(0)
	}

	gzFile.Close()

	updateFiles, err := dirFileList("temp/update/")

	if err != nil {
		fmt.Println("Error with update processing.")
		fmt.Println(err.Error())
		os.Exit(0)
	}

	// Create a backup before continuing
	fmt.Println("Creating backup of old files.")

	var backedUpFiles []string

	for _, file := range updateFiles {
		err = backupFile(updateDir + file)

		if err != nil {
			if err.Error() == "file not found" {
				continue
			}

			fmt.Println("Unable to continue due backup failure.")
			fmt.Println(err.Error())
			os.Exit(0)
		}

		backedUpFiles = append(backedUpFiles, updateDir+file)
	}

	// Update the files
	fmt.Println("Patching files.")

	for _, file := range updateFiles {
		err = copyFile("temp/update/"+file, updateDir+file)

		if err != nil {
			// Failure to update file so restore from backup
			fmt.Println("Failed to patch file: " + file)
			fmt.Println(err.Error())
			fmt.Println("Restoring files from backup")

			for _, backup := range backedUpFiles {
				err = restoreFile(backup)

				if err != nil {
					fmt.Println("Failed to restore file: " + backup)
					fmt.Println(err.Error())
				}
			}

			fmt.Println("Failed to update.")
			os.Exit(0)
		}
	}

	fmt.Println("Update completed.")
	fmt.Println("Cleaning up.")

	err = os.RemoveAll("temp/")

	if err != nil {
		fmt.Println("Error cleaning up, dir \"temp/\" has to be manually deleted.")
		fmt.Println(err.Error())
	}

	os.Exit(1)
}
