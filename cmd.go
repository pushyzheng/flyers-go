package flyers

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
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

func (c *Cmd) Run() (*CmdResult, error) {
	cmd := exec.Command(c.name, c.args...)
	stdout, err := cmd.StdoutPipe()
	stderr, err2 := cmd.StderrPipe()
	if err != nil || err2 != nil {
		return nil, fmt.Errorf("cannot create std pipe: %v", err)
	}
	if err = cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %v", err)
	}
	return &CmdResult{
		cmd:    cmd,
		stdout: stdout,
		stderr: stderr,
	}, nil
}

func (cr *CmdResult) ScanStdout(fn func(line string)) error {
	scanner := bufio.NewScanner(cr.stdout)
	for scanner.Scan() {
		fn(scanner.Text())
	}
	return cr.cmd.Wait()
}
