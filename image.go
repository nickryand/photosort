package main

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

type ImageObj struct {
	FileObj
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
func (i *ImageObj) parseYearMonthDay(tags *exif.Exif) (string, error) {
	dateStr, err := tags.Get(exif.DateTimeOriginal)
	if err != nil {
		return "", err
	}
	dateSlice := strings.Split(dateStr.StringVal(), " ")

	replacer := strings.NewReplacer(":", "")
	return replacer.Replace(dateSlice[0]), nil
}

func (i *ImageObj) GetDestination() (string, error) {
	log.Println("Processing file: ", i.filename)

	// Set the tags on the struct
	tags, err := i.getTags()
	if err != nil {
		return "", err
	}

	date, err := i.parseYearMonthDay(tags)
	if err != nil {
		return "", err
	}

	return path.Join(i.destpath, date), nil
}
