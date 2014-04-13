package main

import (
	"github.com/rwcarlsen/goexif/exif"
	"fmt"
	"flag"
	"os"
	"bufio"
	"io"
	"hash/crc32"
	"crypto/md5"
	"strings"
	"path"
	"sort"
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

func GetCrc32(file *string) (uint32, error) {
	fp, err := os.Open(*file)
	if err != nil {
		return 0, err
	}
	defer fp.Close()

	reader := bufio.NewReader(fp)

	// Calculate the CRC32 using buffered IO
	hasher := crc32.New(crc32.IEEETable)

	_, err = io.Copy(hasher, reader)
	if err != nil {
		return 0, err
	}

	return hasher.Sum32(), nil
}

func GetMd5(file *string) (hash []byte, err error) {
	fp, err := os.Open(*file)
	if err != nil {
		return []byte{}, err
	}
	defer fp.Close()

	reader := bufio.NewReader(fp)

	// Calculate the MD5 using buffered IO
	hasher := md5.New()

	_, err = io.Copy(hasher, reader)
	if err != nil {
		return []byte{}, err
	}

	return hasher.Sum(nil), nil
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

// Check that a directory already exists in our cache. Using a cache to store
// already created directories removes the need to make a system call to check
// that the directory exists every time.
func CheckDir(dest string, dirs []string) []string {
	// Search for dest in our slice of strings
	index := sort.SearchStrings(dirs, dest)
	if index < len(dirs) && dirs[index] == dest {
		// found it
		return dirs
	} else {
		fmt.Println("Directory not in cache:", dest)
		err := os.MkdirAll(dest, 0755)
		if err != nil {
			panic(err)
		}
				// Insert: https://code.google.com/p/go-wiki/wiki/SliceTricks
		// This code uses the append function to grow the slice by 1 element
		// and has the side effect of growing the underlying array if neccessary
		dirs = append(dirs, dest)
		copy(dirs[index+1:], dirs[index:])
		dirs[index] = dest
	}
	return dirs
}

func main() {
	// Process our command line args
	var destpath = flag.String("destpath", ".", "Destination directory for Photos")
	flag.Parse()
	args := flag.Args()

	// Keep a sorted list of directories we have already created
	// so we do not need to make a system call for each picture.
	var dirs []string

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

		// Check to make sure the directory exists. If not create it
		destination := path.Join(*destpath, date)
		dirs = CheckDir(destination, dirs)

		extention := path.Ext(filepath)

		hash, err := GetMd5(&filepath)
		if err != nil {
			panic(err)
		}

		filename := fmt.Sprintf("%X%s", hash, strings.ToLower(extention))
		fmt.Println("Copying file to: ", path.Join(destination, filename))
	}
}
