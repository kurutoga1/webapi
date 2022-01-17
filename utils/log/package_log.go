/*
packageでログを書く場合はこのファイルをコピーし、logfを使用する。
外部パッケージからSetLoggerを使用し、ロガーをセットする。
デフォルトではログを標準出力する。
*/

package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
)

var (
	logger = log.New(os.Stdout, "", log.LstdFlags)
	logMu  sync.Mutex
)

func SetLogger(l *log.Logger) {
	if l == nil {
		l = log.New(os.Stdout, "", log.LstdFlags)
	}
	logMu.Lock()
	logger = l
	logMu.Unlock()
}

func logf(format string, v ...interface{}) {
	// 呼び出されたファイル名、行数を取得
	_, filename, line, _ := runtime.Caller(1)
	baseName := filepath.Base(filename)
	logMu.Lock()
	formatted := fmt.Sprintf(format, v...)
	logStr := baseName + ":" + strconv.Itoa(line) + ": " + formatted
	logger.Println(logStr)
	logMu.Unlock()
}
