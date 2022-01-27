package file

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func FileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}

	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		return true
	}
	return false
}

// CreateSpecifiedFile 特定のサイズのファイルを作成する。
func CreateSpecifiedFile(path string, kb int64) error {
	size := int64(1024 * kb)
	fd, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("CreateSpecifiedFile: %v", err)
	}
	_, err = fd.Seek(size-1, 0)
	if err != nil {
		return fmt.Errorf("CreateSpecifiedFile: %v", err)
	}
	_, err = fd.Write([]byte{0})
	if err != nil {
		return fmt.Errorf("CreateSpecifiedFile: %v", err)
	}
	err = fd.Close()
	if err != nil {
		return fmt.Errorf("CreateSpecifiedFile: %v", err)
	}
	return nil
}

// ReadBytesWithSize はio.ReaderでbufferSizeの制限で読み込み、文字列で返す。
func ReadBytesWithSize(r io.Reader, bufferSize int) (string, error) {
	buffer := make([]byte, bufferSize)
	if _, err := r.Read(buffer); err != nil {
		if err != io.EOF {
			return "", fmt.Errorf("ReadBytesWithSize: %v", err)
		}
	}

	// bufferSizeが満たされていないと[114 111 99 101 115 115 0 0 0 0 0 0 0 0]みたいになるので0を削除する。
	var b []byte
	zero := []byte{0}
	if len(buffer) > 0 && bytes.Contains(buffer, zero) {
		startZero := bytes.IndexByte(buffer, 0)
		b = buffer[:startZero]
	} else {
		b = buffer
	}

	return string(b), nil
}
