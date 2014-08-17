package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/modcloth/go-fileutils"
)

type File interface {
	GetMd5() ([]byte, error)
	GetDestination() (string, error)
	GetFileName() string
	IsCached(string) (bool, int)
	AddToCache(string, int)
	IncDup()
	IncFail()
	IncSuccess()
	getDestpath() string
}

type FileObj struct {
	destpath string
	filename string
	cache    *Cache
}

func (f *FileObj) getDestpath() string {
	return f.destpath
}

func (f *FileObj) IncDup() {
	f.cache.duplicate++
}

func (f *FileObj) IncFail() {
	f.cache.failure++
}

func (f *FileObj) IncSuccess() {
	f.cache.success++
}

func (f *FileObj) IsCached(s string) (bool, int) {
	return f.cache.IsCached(s)
}

func (f *FileObj) AddToCache(s string, index int) {
	f.cache.Insert(s, index)
}

func (f *FileObj) GetFileName() string {
	return f.filename
}

func (f *FileObj) GetMd5() (hash []byte, err error) {
	fp, err := os.Open(f.filename)
	if err != nil {
		log.Printf("Error while attempting to open %s: %s\n", f.filename, err)
		return
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

func (f *FileObj) GetDestination() (string, error) {
	return f.destpath, nil
}

func CheckDir(f File, directory string) error {
	found, index := f.IsCached(directory)
	if !found {
		debug.Printf("Directory not in cache: %s", directory)
		err := os.MkdirAll(directory, 0755)
		if err != nil {
			return err
		}
		f.AddToCache(directory, index)
	}
	return nil
}

func CopyFile(f File) error {
	destdir, err := f.GetDestination()
	if err != nil {
		return err
	}

	// Verify our directory exists
	if err := CheckDir(f, destdir); err != nil {
		return err
	}

	// Reuse the filename extension
	extention := path.Ext(f.GetFileName())

	hash, err := f.GetMd5()
	if err != nil {
		return err
	}

	name := fmt.Sprintf("%X%s", hash, strings.ToLower(extention))
	destination := path.Join(destdir, name)

	// Only copy of the file does not already exist
	if _, err := os.Stat(destination); os.IsNotExist(err) {
		log.Printf("Copying file to: %s", destination)
		if err = fileutils.Cp(f.GetFileName(), destination); err != nil {
			f.IncFail()
			return err
		}
		f.IncSuccess()
	} else {
		f.IncDup()
		debug.Printf("File %s already exists\n", destination)
	}
	return nil
}
