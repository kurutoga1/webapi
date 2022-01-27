/*
ファイルを移動させる機能を提供するパッケージ
*/

package file

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type Mover interface {
	// Move srcからdstにファイルを移動させる
	Move(src, dst string) error
}

func NewMover() Mover {
	return &mover{}
}

type mover struct {
	mu *sync.Mutex
}

// Move ファイルを移動させる
// os.Renameはパーティションを飛び越えてファイルのrenameは出来ないため。
func (m *mover) Move(sourcePath, destPath string) (err error) {
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("MoveFile: %s", err)
	}

	dstFile, err := os.Create(destPath)
	if err != nil {
		err2 := srcFile.Close()
		if err2 != nil {
			return err
		}
		return fmt.Errorf("MoveFile: %s, %s", err, err2)
	}

	defer func(outputFile *os.File) {
		err = outputFile.Close()
	}(dstFile)

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("MoveFile: %s", err)
	}

	err = srcFile.Close()
	if err != nil {
		return fmt.Errorf("MoveFile: %s", err)
	}

	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("MoveFile: %s", err)
	}
	return nil
}
