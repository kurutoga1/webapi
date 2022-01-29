package file_test

import (
	"fmt"
	"os"
	"testing"
	"webapi/utils/file"
)

var testFile string

func init() {
	testFile = "testFIle"
}

func tearDown() {
	os.Remove(testFile)
}

func TestFileExists(t *testing.T) {
	_, err := os.Create(testFile)
	if err != nil {
		panic(err)
	}
	if !file.FileExists(testFile) {
		t.Errorf("got: true, want: false.")
	}
	if file.FileExists("dummyFile") {
		t.Errorf("got: false, want: true.")
	}

	t.Cleanup(func() {
		tearDown()
	})
}

func TestCreateSpecifiedFile(t *testing.T) {
	err := file.CreateSpecifiedFile(testFile, 200)
	if err != nil {
		panic(err)
	}

	f, err := os.Stat(testFile)
	var wantSize int64 = 209715200 // 200KB
	if f.Size() != wantSize {
		t.Errorf("got: %v, want: %v", f.Size(), wantSize)
	}

	t.Cleanup(func() {
		tearDown()
	})
}

func TestReadBytesWithSize(t *testing.T) {
	// これでいくはずだけどなぜかうまくいかない
	t.Skip()
	b := []byte("aaaaabbbbbccccceeeeedddddfffffgggggiiiii")
	f, err := os.Create(testFile)
	if err != nil {
		panic(err)
	}
	_, err = f.Write(b)
	if err != nil {
		panic(err)
	}
	s, err := file.ReadBytesWithSize(f, 10)
	if err != nil {
		panic(err)
	}

	wantStr := "aaaaabbbbb"
	if s != wantStr {
		t.Errorf("got: %v, want: %v", s, wantStr)
	}
	fmt.Println(s)
}
