package main

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

type ImageObj struct {
	*FileObj
}

func (i *ImageObj) getTags() (*exif.Exif, error) {
	fp, err := os.Open(i.filename)
	if err != nil {
		// Issue opening the file
		return nil, err
	}

	tags, err := exif.Decode(fp)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// Parse the date string from an exif tag and return the year, month and day
// in the format YYYYMMDD.
func (i *ImageObj) getDatePath(tags *exif.Exif) (string, error) {
	dateStr, err := tags.Get(exif.DateTimeOriginal)
	if err != nil {
		return "", err
	}
	dateSlice := strings.SplitN(dateStr.StringVal(), ":", 3)

	return path.Join(dateSlice[0], dateSlice[1]), nil
}

func (i *ImageObj) GetDestination() (string, error) {
	log.Println("Processing file: ", i.filename)

	// Set the tags on the struct
	tags, err := i.getTags()
	if err != nil {
		return path.Join(i.destpath, "image"), nil
	}

	date, err := i.getDatePath(tags)
	if err != nil {
		return path.Join(i.destpath, "image"), nil
	}

	return path.Join(i.destpath, date), nil
}
