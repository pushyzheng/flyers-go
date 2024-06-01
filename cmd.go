package flyers

import (
	"bufio"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os/exec"
	"strings"
)

type Cmd struct {
	ctx  context.Context
	name string
	args []string
}

type CmdResult struct {
	cmd    *exec.Cmd
	stdout io.ReadCloser
	stderr io.ReadCloser
}

func NewCmd(ctx context.Context, name string, args ...string) *Cmd {
	return &Cmd{
		ctx:  ctx,
		name: name,
		args: args,
	}
}

func NewShellCmd(ctx context.Context, command string) *Cmd {
	return &Cmd{
		ctx:  ctx,
		name: "sh",
		args: []string{"-c", command},
	}
}

// Run 启动命令, 并返回命令的结果, 不会等待命令的执行完成
func (c *Cmd) Run() (*CmdResult, error) {
	cmd := exec.Command(c.name, c.args...)
	stdout, err := cmd.StdoutPipe()
	stderr, err2 := cmd.StderrPipe()
	if err != nil || err2 != nil {
		return nil, fmt.Errorf("cannot create std pipe: %v", err)
	}
	logrus.Infof("starting cmd: %s %s", c.name, strings.Join(c.args, " "))
	// not blocking
	if err = cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %v", err)
	}
	return &CmdResult{
		cmd:    cmd,
		stdout: stdout,
		stderr: stderr,
	}, nil
}

func (cr *CmdResult) Wait() ([]byte, []byte, error) {
	stdout, err := cr.GetStdout()
	stderr, err2 := cr.GetStderr()
	if err != nil || err2 != nil {
		return nil, nil, fmt.Errorf("stdout/stderr get error: %v %v", err, err2)
	}
	return stdout, stderr, nil
}

// GetStdout 阻塞等待获取标准输出
func (cr *CmdResult) GetStdout() ([]byte, error) {
	return io.ReadAll(cr.stdout)
}

func (cr *CmdResult) GetStdoutStr() string {
	data, err := cr.GetStdout()
	if err != nil {
		return ""
	}
	return string(data)
}

func (cr *CmdResult) GetStderr() ([]byte, error) {
	return io.ReadAll(cr.stderr)
}

func (cr *CmdResult) GetStderrStr() string {
	data, err := cr.GetStderr()
	if err != nil {
		return ""
	}
	return string(data)
}

func (cr *CmdResult) ScanStdout(fn func(line string)) error {
	scanner := bufio.NewScanner(cr.stdout)
	for scanner.Scan() {
		fn(scanner.Text())
	}
	return scanner.Err()
}

func (cr *CmdResult) ScanStderr(fn func(line string)) error {
	scanner := bufio.NewScanner(cr.stderr)
	for scanner.Scan() {
		fn(scanner.Text())
	}
	return scanner.Err()
}
