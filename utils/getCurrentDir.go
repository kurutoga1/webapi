/*
GetCurrentDirを定義したファイル
GetCurrentDirはテストやconfig.jsonを読み込み際に使用する。
*/

package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetFiles ディレクトリ名を渡し、ファイルのリスト、エラーを返す。
func GetFiles(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, file := range files {
		if !file.IsDir() {
			paths = append(paths, filepath.Join(dir, file.Name()))
			continue
		}
	}
	return paths, nil
}

// GetCurrentDir
// コンパイルする前はこの関数を呼び出したファイルがあるディレクトリを返し、
// コンパイルした後はビルドしたファイルを置くディレクトリになる。
func GetCurrentDir() (string, error) {
	// 実際にコマンドを実行しているファイルのカレントディレクトリが入る。
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	currentFiles, err := GetFiles(currentDir)
	if err != nil {
		return "", fmt.Errorf("GetCurrentDir: %v", err)
	}

	// カレントディレクトリに.goファイルがあるかないかでコンパイル前か後か判断している。
	var f bool = false
	for _, file := range currentFiles {
		if strings.Contains(file, ".go") {
			f = true
		}
	}

	if f {
		_, filename, _, _ := runtime.Caller(1)
		return filepath.Dir(filename), nil
	}

	return currentDir, nil
}
