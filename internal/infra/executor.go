package infra

import (
	"bytes"
	"fmt"
	"os/exec"
)

// ExecuteScriptAndGetOutputError 执行脚本并返回 stdout、stderr 和错误
func ExecuteScriptAndGetOutputError(script string) (stdout string, stderr string, err error) {
	cmd := exec.Command("/bin/bash", script)

	var stdoutBuf, stderrBuf bytes.Buffer

	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()
	if err != nil {
		return stdoutBuf.String(), stderrBuf.String(), fmt.Errorf("执行脚本失败: %v\n错误输出:\n%s", err, stderrBuf.String())
	}

	return stdoutBuf.String(), stderrBuf.String(), nil
}
