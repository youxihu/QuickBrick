// infra/script_executor.go

package infra

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// ExecuteScriptAndGetOutputError 执行脚本并返回 stdout、stderr 和错误
// 同时在终端实时输出每行内容
func ExecuteScriptAndGetOutputError(script string) (stdout string, stderr string, err error) {
	cmd := exec.Command("/bin/bash", script)

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	err = cmd.Start()
	if err != nil {
		return "", "", fmt.Errorf("启动脚本失败: %v", err)
	}

	// 实时读取 stdout
	go func() {
		reader := bufio.NewReader(stdoutPipe)
		for {
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				break
			}
			if line != "" {
				fmt.Print(line)             // 实时输出到控制台
				stdoutBuf.WriteString(line) // 同步写入 buffer
			}
			if err == io.EOF {
				break
			}
		}
	}()

	// 实时读取 stderr
	go func() {
		reader := bufio.NewReader(stderrPipe)
		for {
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				break
			}
			if line != "" {
				fmt.Fprint(os.Stderr, line) // 实时输出到 stderr
				stderrBuf.WriteString(line) // 同步写入 buffer
			}
			if err == io.EOF {
				break
			}
		}
	}()

	err = cmd.Wait()
	if err != nil {
		return stdoutBuf.String(), stderrBuf.String(), fmt.Errorf("执行脚本失败: %v\n错误输出:\n%s", err, stderrBuf.String())
	}

	return stdoutBuf.String(), stderrBuf.String(), nil
}
