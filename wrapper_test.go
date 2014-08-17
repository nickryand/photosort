package main

import (
	"reflect"
	"testing"
)

var testCases = []struct {
	filename    string
	destpath    string
	mimetype    string
	wrappedType interface{}
	error       bool
}{
	{"testdata/1_datetime_key.jpg", "/tmp/testtemp", "image/jpeg", &ImageObj{}, false},
	{"testdata/2_datetime_keys.jpg", "/tmp/testtemp", "image/jpeg", &ImageObj{}, false},
	{"testdata/3_datetime_keys.jpg", "/tmp/testtemp", "image/jpeg", &ImageObj{}, false},
	{"testdata/quicktime_sample.mp4", "/tmp/testtemp/video", "video/mp4", &FileObj{}, false},
	{"testdata/test.data", "/tmp/testtemp/other", "binary", &FileObj{}, false},
	{"testdata/test.png", "/tmp/testtemp/image", "image/png", &FileObj{}, false},
	// This ensures that proper error handling exists for arguments that
	// do not actually exist.
	{"testdata/doesnt-exist", "", "", nil, true},
}

func TestGetWrappedFile(t *testing.T) {
	for _, tc := range testCases {
		fobj := FileObj{
			filename: tc.filename,
			destpath: "/tmp/testtemp",
		}
		wrapped, err := getWrappedFile(&fobj)
		if err != nil {
			if !tc.error {
				t.Errorf("Unexpected error wrapping object: %s", err)
			}
			continue
		}

		type1 := reflect.TypeOf(tc.wrappedType)
		type2 := reflect.TypeOf(wrapped)
		if type1 != type2 {
			t.Errorf("Wrapped type should be %s got %s",
				reflect.TypeOf(tc.wrappedType), reflect.TypeOf(wrapped))
		}

		fileObj := wrapped.(File)
		if tc.destpath != fileObj.getDestpath() {
			t.Errorf("Despath should be %s got %s",
				tc.destpath, fileObj.getDestpath())
		}
	}
}

func TestGetFileType(t *testing.T) {
	for _, tc := range testCases {
		filetype, err := getFileType(tc.filename)
		if err != nil {
			if !tc.error {
				t.Errorf("Error getting type: %s", err)
			}
			continue
		}

		if filetype != tc.mimetype {
			t.Errorf("Incorrect filetype reported, should be %s got %s", tc.mimetype, filetype)
		}
	}
}
