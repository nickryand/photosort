package main

import (
	"path"
	"strings"
	"os"

	"github.com/nickryand/magic"
)

// Wrap a FileObj if there is support for the File's MIME type.
// Otherwise, use the MIME type determine the destination path
// basename.
func getWrappedFile(file *FileObj) (interface{}, error) {
	mimeType, err := getFileType(file.filename)
	if err != nil {
		return nil, err
	}

	switch {
	case strings.HasPrefix(mimeType, "image/jpeg"):
		return &ImageObj{file}, nil
	case strings.HasPrefix(mimeType, "image"):
		file.destpath = path.Join(file.destpath, "image")
	case strings.HasPrefix(mimeType, "video"):
		file.destpath = path.Join(file.destpath, "video")
	default:
		file.destpath = path.Join(file.destpath, "other")
	}
	return file, nil
}

func getFileType(path string) (string, error) {
	// Ensure file actually exists before we attempt to figure out what
	// type of file it is.
	if _, err := os.Stat(path); err != nil {
		return "", err
	}

	conn, err := magic.Open(magic.FlagMimeType)
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
