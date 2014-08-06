package main

import (
	"strings"

	"github.com/nickryand/magic"
)

func getWrappedFile(file *FileObj) (File, error) {
	output, err := getFileType(file.filename)
	if err != nil {
		return nil, err
	}

	switch {
	case strings.HasPrefix(output, "image"):
		debug.Printf("Filetype is image\n")
		return &ImageObj{file}, nil
	}
	return file, nil
}

func getFileType(path string) (string, error) {
	conn, err := magic.Open(magic.FlagMime)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	if err = conn.Load(""); err != nil {
		return "", err
	}

	output, err := conn.File(path)
	if err != nil {
		return "", err
	}

	return output, nil
}
