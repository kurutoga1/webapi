package kernel

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// SimpleExec は実行するためのコマンドをもらい、実行し、stdout, stderr, errを返す
func SimpleExec(command string) (stdoutStr string, stderrStr string, err error) {

	cmd := GetCmdFromStr(command)

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	stdoutStr = stdout.String()
	stderrStr = stderr.String()

	return
}

// GetCmdFromStr コマンド文字列を受け取り、*exec.Cmdを返す
func GetCmdFromStr(command string) *exec.Cmd {
	commands := strings.Split(command, " ")
	return exec.Command(commands[0], commands[1:]...)
}

type ExecuteTimeOutError struct {
	msg string
}

func (e *ExecuteTimeOutError) Error() string {
	return fmt.Sprintf("execute timeout: " + e.msg)
}

// ExecuteWithTimeout は実行するためのコマンド, 時間制限をもらい、OutputManagerインタフェースを返す
func ExecuteWithTimeout(command string, timeOut int) (stdout bytes.Buffer, stderr bytes.Buffer, err error) {
	fName := "ExecuteWithTimeout"

	cmd := GetCmdFromStr(command)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Start()
	if err != nil {
		return stdout, stderr, fmt.Errorf("%v: %v", fName, err)
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
		err = cmd.Process.Kill()
		if err != nil {
			return stdout, stderr, fmt.Errorf("%v: %v", fName, &ExecuteTimeOutError{err.Error()})
		}
		return stdout, stderr, fmt.Errorf("%v: %v", fName, &ExecuteTimeOutError{})

	case err = <-done:
		if err != nil {
			return stdout, stderr, fmt.Errorf("%v: %v", fName, err)
		}
		return
	}
}
