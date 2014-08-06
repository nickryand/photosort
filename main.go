package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

func ProcessDirectory(dirname, destpath, other string, cache *Cache) {
	walkFn := func(filename string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			ProcessFile(filename, destpath, other, cache)
		}
		return nil
	}
	filepath.Walk(dirname, walkFn)
}

func ProcessFile(filename, destpath, other string, cache *Cache) {
	file := FileObj{
		filename: filename,
		destpath: destpath,
		other:    other,
		cache:    cache,
	}

	wrapped, err := getWrappedFile(&file)
	if err != nil {
		log.Printf("Error while wrapping a file: %s", err)
		return
	}

	if err := CopyFile(wrapped); err != nil {
		log.Printf("Error trying to copy file: %s", err)
	}
}

func main() {
	// Process our command line args
	destpath := flag.String("destpath", ".", "Destination directory for Sorted Photos")
	otherdir := flag.String("other", "other", "Name of directory to store non-image files")

	flag.Parse()
	args := flag.Args()

	// Keep a sorted list of directories we have already created
	// so we do not need to make a system call for each picture.
	var cache Cache
	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			log.Printf("Error while attempting to stat %s: %s", arg, err)
		}

		if info.IsDir() {
			ProcessDirectory(arg, *destpath, *otherdir, &cache)
		} else {
			ProcessFile(arg, *destpath, *otherdir, &cache)
		}
	}

	log.Printf("================== Report ==================")
	log.Printf("Duplicates: %d\n", cache.duplicate)
	log.Printf("Successes: %d\n", cache.success)
	log.Printf("Failures: %d\n", cache.failure)
}
