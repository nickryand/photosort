package main

import (
	"github.com/rwcarlsen/goexif/exif"
	"fmt"
	"flag"
	"os"
	"io"
	"hash/crc32"
	"strings"
	"path"
)

func GetTags(file *string) (*exif.Exif, error) {
	fp, err := os.Open(*file)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	tags, err := exif.Decode(fp)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// Return the CRC32 of a file on disk.
func GetCrc32(file *string) (uint32, error) {
	var checksum uint32
	var count, total int

	fp, err := os.Open(*file)
	if err != nil {
		return checksum, err
	}
	defer fp.Close()

	data := make([]byte, 8192)
	count, err = fp.Read(data)
	for {
		if err == io.EOF {
			break
		} else if err != nil {
			return checksum, err
		}

		checksum = crc32.Update(checksum, crc32.IEEETable, data[:count])

		count, err = fp.Read(data)
		total += count
	}

	return checksum, nil
}

// Parse the date string from an exif tag and return the year, month and day
// in the format YYYYMMDD.
func ParseYearMonthDay(tags *exif.Exif) (string, error) {
	dateStr, err := tags.Get(exif.DateTimeOriginal)
	if err != nil {
		return "", err
	}
	dateSlice := strings.Split(dateStr.StringVal(), " ")

	replacer := strings.NewReplacer(":", "")
	return replacer.Replace(dateSlice[0]), nil
}

func main() {
	var destpath = flag.String("destpath", ".", "Destination directory for Photos")
	flag.Parse()
	args := flag.Args()

	for _, filepath := range args {
		fmt.Println("Processing file: ", filepath)

		tags, err := GetTags(&filepath)
		if err != nil {
			panic(err)
		}

		date, err := ParseYearMonthDay(tags)
		if err != nil {
			panic(err)
		}

		destination := path.Join(*destpath, date)
		extention := path.Ext(filepath)

		crc, err := GetCrc32(&filepath)
		if err != nil {
			panic(err)
		}

		filename := fmt.Sprintf("%X%s", crc, strings.ToLower(extention))
		fmt.Println("Copying file to: ", path.Join(destination, filename))
	}
}
