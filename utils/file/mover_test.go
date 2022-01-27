package file_test

import (
	"os"
	"testing"
	"webapi/utils/file"
)

func TestMove(t *testing.T) {
	b := []byte("aaaaabbbbbccccceeeeedddddfffffgggggiiiii")
	f, err := os.Create(testFile)
	if err != nil {
		panic(err)
	}
	_, err = f.Write(b)
	if err != nil {
		panic(err)
	}
	s, err := f.Stat()
	if err != nil {
		panic(err)
	}
	dstFile := "dst"
	mover := file.NewMover()
	err = mover.Move(testFile, dstFile)
	if err != nil {
		panic(err)
	}

	if !file.FileExists(dstFile) {
		t.Errorf("got: false, want: true")
	}

	d, err := os.Stat(dstFile)
	if err != nil {
		panic(err)
	}

	if s.Size() != d.Size() {
		t.Errorf("got: %v, want: %v", s.Size(), d.Size())
	}

	os.Remove(dstFile)

}
