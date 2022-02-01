/*
executer.go
ContextManagerの値をベースにコマンドを実行する。
コマンドを実行した標準出力、出力ファイルなどはOutputManagerにセットする。
実行した後に使用したディレクトリは一定時間後、削除する。
*/

package executer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"webapi/pro/execution/contextManager"
	"webapi/pro/execution/outputManager"
	"webapi/pro/msgs"
	"webapi/utils/execution"
	"webapi/utils/http"
)

// Executer はコマンドを実行する構造体のインタフェース
// 返り値はOutputManager(インターフェース) エラーもoutputManagerの中に入れる、
type Executer interface {
	Execute(contextManager.ContextManager) outputManager.OutputManager
}

// NewExecuter はfileExecuter構造体を返す。
func NewExecuter() Executer {
	return &fileExecuter{}
}

type fileExecuter struct{}

// errorOutWrap は中の３行は頻繁に使用するので行数削減と見やすくするため
// OutputManagerの中にセットする
func errorOutWrap(out outputManager.OutputManager, err error, status string) outputManager.OutputManager {
	out.SetErrorMsg(err.Error())
	out.SetStatus(status)
	return out
}

func (f *fileExecuter) Execute(ctx contextManager.ContextManager) (out outputManager.OutputManager) {

	// コマンド実行
	out = Exec(ctx.Command(), ctx.Config().ProgramTimeOut, ctx.Config().StdoutBufferSize, ctx.Config().StderrBufferSize)
	if out.Status() != msgs.OK {
		return
	}

	// 出力ファイルたちはまだ通常のパスなのでそれを
	// CURLで取得するためにURLパスに変換する。
	outFileURLs, err := GetOutFileURLs(ctx.OutputDir(), ctx.Config().ServerIP, ctx.Config().ServerPort, ctx.Config().FileServer.Dir)
	if err != nil {
		out.SetStatus(msgs.SERVERERROR)
		out.SetErrorMsg(err.Error())
		return
	}

	// 時間経過後ファイルを削除
	go func() {
		err := DeleteDirSomeTimeLater(ctx.ProgramTempDir(), ctx.Config().DeleteProcessedFileLimitSecondTime)
		if err != nil {
			fmt.Printf("Execute: %v \n", err)
		}
	}()

	out.SetOutURLs(outFileURLs)

	return
}

// Exec は実行するためのコマンド, 時間制限をもらい、OutputManagerインタフェースを返す
// エラーがでた場合もoutputInfoのエラーメッセージの中に格納する。
func Exec(command string, timeOut int, stdOutBufferSize, stdErrBufferSize int) (out outputManager.OutputManager) {
	out = outputManager.NewOutputManager()

	var timeoutError *execution.TimeoutError
	stdout, stderr, err := execution.ExecuteWithTimeout(command, timeOut)

	if err1 := out.SetStdOut(&stdout, stdOutBufferSize); err1 != nil {
		out.SetStatus(msgs.SERVERERROR)
		out.SetErrorMsg(fmt.Sprintf("err: %v, err1: %v", err.Error(), err1.Error()))
		return
	}

	if err2 := out.SetStdErr(&stderr, stdErrBufferSize); err2 != nil {
		out.SetStatus(msgs.SERVERERROR)
		out.SetErrorMsg(fmt.Sprintf("err: %v, err2: %v", err.Error(), err2.Error()))
		return
	}

	if err != nil {
		if errors.As(err, &timeoutError) {
			// プログラムがタイムアウトした場合
			out.SetStatus(msgs.PROGRAMTIMEOUT)
			out.SetErrorMsg(err.Error())
			return
		} else {
			// プログラムがエラーで終了した場合
			out.SetStatus(msgs.PROGRAMERROR)
			out.SetErrorMsg(err.Error())
		}
	} else {
		// 正常終了した場合
		out.SetStatus(msgs.OK)
	}

	return
}

// DeleteDirSomeTimeLater は一定時間後にディレクトリを削除する
func DeleteDirSomeTimeLater(dirPath string, seconds int) error {
	// wait some seconds
	time.Sleep(time.Second * time.Duration(seconds))
	err := os.RemoveAll(dirPath)
	if err != nil {
		return fmt.Errorf("DeleteDirSomeTimeLater: %v", err)
	}
	return nil
}

// GetOutFileURLs はコマンドを実行した後に使用する。
// プログラム出力ディレクトリの全てのファイルを取得するURLのリストを返す。
func GetOutFileURLs(outputDir string, serverIP, serverPort, fileServerDir string) ([]string, error) {
	// 出力されたディレクトリの複数ファイルをglobで取得
	pattern := outputDir + "/*"
	outFiles, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("GetOutFileURLs: %v ", err)
	}

	outFileURLs := make([]string, 0, 20)
	for _, outfile := range outFiles {
		outFileURL, err := http.GetURLFromFilePath(outfile, serverIP, serverPort, fileServerDir)
		if err != nil {
			return nil, fmt.Errorf("GetOutFileURLs: %v", err)
		}
		// サーバがwindowsだった場合、出力パス区切りを¥から/に変更する。
		outFileURL = strings.Replace(outFileURL, "¥", "/", -1)
		outFileURLs = append(outFileURLs, outFileURL)
	}

	return outFileURLs, nil
}
