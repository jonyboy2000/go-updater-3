package main

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Gets the file list of a given directory and any subdirectories
// Returns an error (should be nil) and string Slice of files
func dirFileList(dir string) ([]string, error) {
	directory, err := os.Stat(dir)

	if err != nil {
		return nil, err
	}

	if !directory.IsDir() {
		return nil, errors.New("The 'directory' to list files for is not a directory: " + directory.Name())
	}

	directoryList, err := ioutil.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	var filenames []string

	for _, file := range directoryList {
		if file.IsDir() {
			// Recursively call fileNames
			files, err := dirFileList(dir + "/" + file.Name())

			for _, subDirFile := range files {
				filenames = append(filenames, file.Name()+"/"+subDirFile)
			}

			if err != nil {
				return nil, err
			}

		} else {
			filenames = append(filenames, file.Name())
		}
	}

	return filenames, nil
}

// Backs up the file by duplicating it and appending .bak to the extension
// If the file doesn't exist, it will return an error message of "file not
// found", which should be handled as required.
func backupFile(file string) error {
	// Checks source file exists
	_, err := os.Stat(file)

	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("file not found")
		}
		return err
	}

	err = copyFile(file, file+".bak")

	if err != nil {
		return err
	}

	return nil
}

// Restores the file by copying it from the .bak duplicate
// Please note that the .bak extension should not be passed to the function
func restoreFile(file string) error {
	// Checks source file exists
	_, err := os.Stat(file + ".bak")

	if err != nil {
		return err
	}

	err = copyFile(file+".bak", file)

	if err != nil {
		return err
	}

	return nil
}

// Copies file from src to dest
// Returns a non-nil error if failure
func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)

	if err != nil {
		return err
	}

	dir, _ := filepath.Split(dest)

	if _, err = os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dir, 0755)
		} else {
			return err
		}
	}

	destFile, err := os.Create(dest)

	if err != nil {
		return err
	}

	_, err = io.Copy(destFile, srcFile)

	if err != nil {
		return err
	}

	return nil
}
