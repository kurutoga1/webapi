/*
サーバからファイルをダウンロードする機能を提供するパッケージ
*/

package download

import (
	"path/filepath"
	"sync"
	"webapi/utils/file"
	"webapi/utils/kernel"
)

var (
	currentDir string
	uploadFile string
)

type Downloader interface {
	// Download はダウンロードしたいファイルURLを入れて、outputDirへダウンロードする。
	Download(url string, outputDir string, done chan error, wg *sync.WaitGroup, mover file.Mover)
}

func NewDownloader() Downloader {
	return &downloader{}
}

type downloader struct{}

// Download はダウンロードしたいファイルURLを入れて、outputDirへダウンロードする。
func (d *downloader) Download(url, outputDir string, done chan error, wg *sync.WaitGroup, mover file.Mover) {
	defer wg.Done() // 関数終了時にデクリメント
	command := "curl -OL " + url
	_, _, err := kernel.Exec(command)
	if err != nil {
		done <- err
		return
	}

	// 引数で指定された出力ディレクトリに移動させる
	basename := filepath.Base(url)
	newLocation := filepath.Join(outputDir, basename)
	err = mover.Move(basename, newLocation)
	if err != nil {
		done <- err
		return
	}
	done <- nil
	return
}
