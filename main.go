package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
)

// Version of updater can be found by calling the executable on
// its own with just "--version"
var version = "1.1"

// Main usage:
//  Program should be called with the following arguments:
//  - Update download URL to a .tar.gz archive
//  - Update folder target to overwrite, trailing slash required

func main() {
	// Checks arguments
	argCount := len(os.Args)

	// Checks for version argument
	if argCount == 2 {
		if os.Args[1] == "--version" {
			fmt.Println("Go-Updater Version " + version)
			os.Exit(0)
		}
	}

	// Check arguments
	if argCount < 3 {
		fmt.Println("Invalid arguments provided.")
		os.Exit(1)
	}

	// URL Checks
	var downloadURL string // Declared here to allow building with "go build"
	_, err := url.ParseRequestURI(os.Args[1])
	if err != nil {
		fmt.Println("Invalid URL argument provided.")
		fmt.Println(err.Error())
		os.Exit(1)
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
		os.Exit(1)
	}

	if !updateTarget.IsDir() {
		fmt.Println("The update target dir appears to be a file.")
		os.Exit(1)
	}

	updateDir := os.Args[2]

	autoStart := false
	var autoStartApplication string

	// Checks remaining arguments using regexes
	if argCount != 3 {
		for i := 3; i < argCount-1; i++ {
			if match, _ := regexp.MatchString("(^--start)", os.Args[i]); match {
				autoStart = true
				autoStartApplication = os.Args[i+1]
			}
		}
	}

	// Downloads the .gz to a temporary file
	err = os.MkdirAll("temp/", 0755)

	if err != nil {
		fmt.Println("Unable to create temporary directory.")
		fmt.Println(err.Error())
		cleanUp()
		os.Exit(1)
	}

	archive, err := os.Create("temp/" + path.Base(downloadURL))

	if err != nil {
		fmt.Println("Can not create the temporary download file.")
		fmt.Println(err.Error())
		cleanUp()
		os.Exit(1)
	}

	resp, err := http.Get(downloadURL)

	if err != nil {
		fmt.Println("Can not download the update.")
		fmt.Println(err.Error())
		cleanUp()
		os.Exit(1)
	}

	defer resp.Body.Close()

	n, err := io.Copy(archive, resp.Body)
	if err != nil {
		fmt.Println("Can not download the update.")
		fmt.Println(err.Error())
		cleanUp()
		os.Exit(1)
	}

	fmt.Println("Successfully downloaded the update: " + strconv.Itoa(int(n)) + "bytes downloaded.")
	fmt.Println("Extracting update archive.")

	err = extractArchive("temp/"+path.Base(downloadURL), "temp/update/")

	if err != nil {
		fmt.Println("Error extracting archive.")
		fmt.Println(err.Error())
		cleanUp()
		os.Exit(1)
	}

	archive.Close()

	updateFiles, err := dirFileList("temp/update/")

	if err != nil {
		fmt.Println("Error with update processing.")
		fmt.Println(err.Error())
		cleanUp()
		os.Exit(1)
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
			cleanUp()
			os.Exit(1)
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
			cleanUp()
			os.Exit(1)
		}
	}

	fmt.Println("Update completed.")
	cleanUp()

	// Start application if applicable
	if autoStart {
		fmt.Println("Starting application.")
		application := exec.Command(autoStartApplication)

		err := application.Start()

		if err != nil {
			fmt.Println("Failed to start application")
			fmt.Println(err.Error())
		}
	}

	os.Exit(0)
}

// Clean up function
func cleanUp() {
	fmt.Println("Cleaning up.")

	err := os.RemoveAll("temp/")

	if err != nil {
		fmt.Println("Error cleaning up, dir \"temp/\" has to be manually deleted.")
		fmt.Println(err.Error())
	}
}
