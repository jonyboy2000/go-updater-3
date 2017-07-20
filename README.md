# go-updater
A simple go application for updating files/applications

- Uses "github.com/mholt/archiver" for archive decompression 
  - Currently only .tar.gz is supported although more support can be addded if requested
- Automatically backs up files and will restore if an error occurs during update process

## General Usage
- Downloads update from a sample URL and installs to the "application" dir -> Will overwrite all existing files
  - `go-updater https://example.net/test.tar.gz application/`

- Optional flags can be added like below:
 - `go-updater https://example.net/test.tar.gz application/ --start application.exe`
 - See below for more detailed information regarding the optional flags

## Optional flags:
  - `--start application.exe`
    - Starts the specified application after successful update

## Gets go-updater version version:
  - `go-updater --version`
  - Exits after printing the version string
