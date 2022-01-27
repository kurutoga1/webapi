package execution_test

import (
	"errors"
	"strings"
	"testing"
	"webapi/utils/execution"
)

func TestSimpleExec(t *testing.T) {
	word := "hello world"
	cmd := "echo " + word
	stdout, stderr, err := execution.SimpleExec(cmd)

	if err != nil {
		t.Errorf("err: %v", err)
	}

	// 改行が含まれるため
	if !strings.Contains(stdout, word) {
		t.Errorf("got: %v, want: %v", stdout, word)
	}

	if stderr != "" {
		t.Errorf("got: %v, want: empty string", stderr)
	}
}

func TestExecuteWithTimeout(t *testing.T) {
	var timeoutError *execution.TimeoutError

	cmd := "sleep 5"
	_, _, err := execution.ExecuteWithTimeout(cmd, 1)

	if err != nil {
		if !errors.As(err, &timeoutError) {
			t.Errorf("got: %v, want: %v", err, timeoutError)
		}
	}

	cmd = "sleep 1"
	_, _, err = execution.ExecuteWithTimeout(cmd, 2)

	if err != nil {
		t.Errorf("got: %v, want: not error", err)
	}

	// wrong command
	// この場合はtimeoutErrorではなく通常のexecが出すエラー
	cmd = "seep 1"
	_, _, err = execution.ExecuteWithTimeout(cmd, 2)

	if err != nil && errors.As(err, &timeoutError) {
		t.Errorf("got: %v, want: exec: \"seep\": executable file not found in $PATH error", err)
	}
}
