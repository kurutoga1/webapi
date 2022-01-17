package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func getFiles(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if !file.IsDir() {
			paths = append(paths, filepath.Join(dir, file.Name()))
			continue
		}
	}
	return paths
}

func getCurrentDir() (string, error) {
	// コンパイルする前とコンパイルした後ではファイル構造が異なる。
	// するとtestする時に、場所が変更になり困るので、どっちもに対応したカレントディレクトリを
	// 取得する関数。
	// コンパイルする前はconfディレクトリになり、
	// コンパイルした後はビルドしたファイルを置くディレクトリになる。

	// 実際にコマンドを実行しているファイルのカレントディレクトリが入る。
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	currentFiles := getFiles(currentDir)

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
