package kernel

import (
	"bytes"
	"os/exec"
	"strings"
)

// Exec は実行するためのコマンドをもらい、実行し、stdout, stderr, errを返す
func Exec(command string) (stdoutStr string, stderrStr string, cmderr error) {

	commands := strings.Split(command, " ")

	cmd := exec.Command(commands[0], commands[1:]...)

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	stdoutStr = stdout.String()
	stderrStr = stderr.String()
	cmderr = err

	return
}
