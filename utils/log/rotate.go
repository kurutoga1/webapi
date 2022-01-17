package log

import (
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"webapi/utils"
)

type Rotater interface {
	// Rotate はログのファイルパスとサイズを渡すことにより、ファイルサイズなるようにファイルを上から削る。
	// maxSizeを超えていたら、shavingSizeまでファイルを削る
	Rotate() error
}

func NewLogRotater(shavingSize, maxSize int, mu *sync.Mutex, logger *log.Logger, logFile string) *logRotater {
	return &logRotater{
		shavingSize: shavingSize,
		maxSize:     maxSize,
		mu:          mu,
		logger:      logger,
		logFile:     logFile,
	}
}

type logRotater struct {
	shavingSize int
	maxSize     int
	mu          *sync.Mutex
	logger      *log.Logger
	logFile     string
}

// Rotate はログのファイルパスとサイズを渡すことにより、ファイルサイズなるようにファイルを上から削る。
// maxSizeを超えていたら、shavingSizeまでファイルを削る
func (l *logRotater) Rotate() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !utils.FileExists(l.logFile) {
		return errors.New(l.logFile + "is not found.")
	}

	r, err := os.Open(l.logFile)
	if err != nil {
		return err
	}
	fileInfo, err := r.Stat()
	if err != nil {
		return err
	}
	// ログファイルのサイズがマックスより小さかったら何もしない
	fileSize := fileInfo.Size()
	if fileSize < int64(l.maxSize) {
		return nil
	}

	// 後で書き込むbytes
	var lines []byte
	var seek int64 = 0
	readOK := false
	// ファイルサイズから削るサイズを引き、読み始めるべき場所を格納する
	readOKSeek := fileSize - int64(l.shavingSize)
	scanner := bufio.NewScanner(r)
	// 一行ずつ取り出す
	for scanner.Scan() {
		// seekが読み始めるべき場所を超えたら読み込みフラグをtrueにする
		if readOKSeek < seek {
			readOK = true
		}
		if readOK {
			lines = append(lines, scanner.Bytes()...)
			lines = append(lines, []byte("\n")...)
		}
		// 読み込んだ分seekを進める
		seek += int64(len(scanner.Bytes()))
	}
	err = r.Close()
	if err != nil {
		return err
	}

	// tmpファイル作成
	tempFile, err := ioutil.TempFile(os.TempDir(), "tmp_log.txt")
	if err != nil {
		return err
	}
	w, err := os.Create(tempFile.Name())
	if err != nil {
		return err
	}
	// tmpファイルに書き込み
	_, err = w.Write(lines)
	if err != nil {
		return err
	}
	// tmpファイルに書き込みが成功したので元のログファイルは削除する
	err = os.Remove(l.logFile)
	if err != nil {
		return err
	}
	// tmpファイルを元のログファイル名に変更する
	err = os.Rename(tempFile.Name(), l.logFile)
	if err != nil {
		return err
	}

	// ファイルを変えたのでファイルロガーも変更する必要がある。
	logger2 := GetLogger(l.logFile)
	logger = logger2

	return nil
}
