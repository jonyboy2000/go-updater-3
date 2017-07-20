package main

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mholt/archiver"
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
	defer srcFile.Close()

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
	defer destFile.Close()

	if err != nil {
		return err
	}

	_, err = io.Copy(destFile, srcFile)

	if err != nil {
		return err
	}

	return nil
}

// Extracts the specified archive to the specified directory
// Supports the archive types supported by github.com/mholt/archiver
// Uses the file extension to determine the necessary archive
func extractArchive(archive, dest string) error {
	// Checks archive exists
	_, err := os.Stat(archive)

	if err != nil {
		return err
	}

	// Detects archive type by using the file extension
	// Uses filepath.Ext to start with, then string manipulation
	switch filepath.Ext(archive) {
	case ".zip":
		err = archiver.Zip.Open(archive, dest)
	case ".tar":
		err = archiver.Tar.Open(archive, dest)
	case ".tgz":
		err = archiver.TarGz.Open(archive, dest)
	case ".tbz2":
		err = archiver.TarBz2.Open(archive, dest)
	case ".txz":
		err = archiver.TarXZ.Open(archive, dest)
	case ".tlz4":
		err = archiver.TarLz4.Open(archive, dest)
	case ".tsz":
		err = archiver.TarSz.Open(archive, dest)
	case ".rar":
		err = archiver.Rar.Open(archive, dest)
	default:
		switch archive[len(archive)-8:] {
		case ".tar.bz2":
			err = archiver.TarBz2.Open(archive, dest)
		case ".tar.lz4":
			err = archiver.TarLz4.Open(archive, dest)
		default:
			switch archive[len(archive)-7:] {
			case ".tar.gz":
				err = archiver.TarGz.Open(archive, dest)
			case ".tar.xz":
				err = archiver.TarXZ.Open(archive, dest)
			case ".tar.sz":
				err = archiver.TarSz.Open(archive, dest)
			default:
				err = errors.New("Archive type is not supported.")
			}
		}
	}

	return err
}
