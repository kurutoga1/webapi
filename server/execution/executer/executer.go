/*
executer.go
ContextManagerの値をベースにコマンドを実行する。
コマンドを実行した標準出力、出力ファイルなどはOutputManagerにセットする。
実行した後に使用したディレクトリは一定時間後、削除する。
*/

package executer

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"webapi/server/execution/contextManager"
	"webapi/server/execution/msgs"
	"webapi/server/outputManager"
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
	out = outputManager.NewOutputManager()

	// コマンド実行
	out, err := Exec(ctx.Command(), ctx.Config().ProgramTimeOut, ctx.Config().StdoutBufferSize, ctx.Config().StderrBufferSize)
	if err != nil {
		return errorOutWrap(out, err, msgs.PROGRAMERROR)
	}

	// 出力ファイルたちはまだ通常のパスなのでそれを
	// CURLで取得するためにURLパスに変換する。
	outFileURLs, err := GetOutFileURLs(ctx.OutputDir(), ctx.Config().ServerIP, ctx.Config().ServerPort, ctx.Config().FileServer.Dir)
	if err != nil {
		return errorOutWrap(out, err, msgs.SERVERERROR)
	}

	// 時間経過後ファイルを削除
	go func() {
		err := DeleteDirSomeTimeLater(ctx.ProgramTempDir(), ctx.Config().DeleteProcessedFileLimitSecondTime)
		if err != nil {
			fmt.Printf("Execute: %v \n", err)
		}
	}()

	out.SetOutURLs(outFileURLs)

	out.SetStatus(msgs.OK)
	return
}

// Exec は実行するためのコマンド, 時間制限をもらい、OutputManagerインタフェースを返す
func Exec(command string, timeOut int, stdOutBufferSize, stdErrBufferSize int) (outputManager.OutputManager, error) {
	var outputInfo = outputManager.NewOutputManager()

	commands := strings.Split(command, " ")

	cmd := exec.Command(commands[0], commands[1:]...)

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		return outputInfo, fmt.Errorf("Exec: %v ", err)
	}

	// cmd.Wait()はコマンドの中でエラーが出たらerrorをdoneチャネルに送信。エラーがない場合はnilを送る。
	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	// タイマーが終了したらチャネルを受け取るtimeoutチャネルを定義
	timeout := time.After(time.Second * time.Duration(timeOut))

	// timeoutチャネルが先に来た場合はコマンド実行のプロセスを終了する。
	// doneチャネルが来た場合はどっちにしろコマンドが終了したという合図。
	select {

	case <-timeout:
		// プログラムが設定した時間以内に終了せずタイムアウトする場合。
		err := cmd.Process.Kill()
		if err != nil {
			return outputInfo, fmt.Errorf("Exec: %v ", err)
		}
		outputInfo.SetStatus(msgs.PROGRAMTIMEOUT)
		//if err := outputInfo.SetStdOut(&stdout, cfg.StdoutBufferSize); err != nil {
		if err := outputInfo.SetStdOut(&stdout, stdOutBufferSize); err != nil {
			return outputInfo, fmt.Errorf("Exec: %v ", err)
		}

		//if err := outputInfo.SetStdErr(&stderr, cfg.StderrBufferSize); err != nil {
		if err := outputInfo.SetStdErr(&stderr, stdErrBufferSize); err != nil {
			return outputInfo, fmt.Errorf("Exec: %v ", err)
		}
		return outputInfo, errors.New("program time out error")

	case err := <-done:

		// プログラムが設定した時間以内に終了した場合
		outputInfo.SetStatus(msgs.OK)
		//if err := outputInfo.SetStdOut(&stdout, cfg.StdoutBufferSize); err != nil {
		if err := outputInfo.SetStdOut(&stdout, stdOutBufferSize); err != nil {
			return outputInfo, fmt.Errorf("Exec: %v ", err)
		}

		//if err := outputInfo.SetStdErr(&stderr, cfg.StderrBufferSize); err != nil {
		if err := outputInfo.SetStdErr(&stderr, stdErrBufferSize); err != nil {
			return outputInfo, fmt.Errorf("Exec: %v ", err)
		}
		if err != nil {
			outputInfo.SetStatus(msgs.PROGRAMERROR)
			return outputInfo, fmt.Errorf("Exec: %v ", err)
		}
	}

	return outputInfo, nil
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
