package kernel_test

import (
	"strings"
	"testing"
	"webapi/utils/kernel"
)

func TestExec(t *testing.T) {
	word := "hello world"
	cmd := "echo " + word
	stdout, stderr, err := kernel.SimpleExec(cmd)

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
