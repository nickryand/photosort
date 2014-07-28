package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

func ProcessDirectory(dirname string, destpath string, cache *Cache) {
	// Defined as anonymous for scoping reasons
	walkFn := func(filename string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			ProcessFile(filename, destpath, cache)
		}
		return nil
	}
	filepath.Walk(dirname, walkFn)
}

func ProcessFile(filename, destpath string, cache *Cache) {
	file := ImageObj{
		FileObj{
			filename: filename,
			destpath: destpath,
			cache:    cache,
		},
	}
	if err := CopyFile(&file); err != nil {
		log.Printf("Error trying to copy file: %s", err)
	}
}

func main() {
	// Process our command line args
	destpath := flag.String("destpath", ".", "Destination directory for Sorted Photos")
	//otherdir := flag.String("other", "other", "Name of directory to store non-image files")

	flag.Parse()
	args := flag.Args()

	// Keep a sorted list of directories we have already created
	// so we do not need to make a system call for each picture.
	var dircache Cache
	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			log.Printf("Error while attempting to stat %s: %s", arg, err)
		}

		if info.IsDir() {
			ProcessDirectory(arg, *destpath, &dircache)
		} else {
			ProcessFile(arg, *destpath, &dircache)
		}
	}
}
