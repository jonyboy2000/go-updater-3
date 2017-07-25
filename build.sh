#!/bin/bash
GOOS=linux   GOARCH=386   go build -o releases/linux_386/go-update
GOOS=linux   GOARCH=amd64 go build -o releases/linux_amd64/go-update
GOOS=linux   GOARCH=arm   go build -o releases/linux_arm7/go-update
GOOS=linux   GOARCH=arm64 go build -o releases/linux_arm64/go-update

GOOS=darwin  GOARCH=amd64 go build -o releases/mac_amd64/go-update

GOOS=windows GOARCH=386   go build -o releases/windows_386/go-update.exe
GOOS=windows GOARCH=amd64 go build -o releases/windows_amd64/go-update.exe
