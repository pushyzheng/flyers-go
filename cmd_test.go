package flyers

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNewShellCmd(t *testing.T) {
	t.Run("GetStdout", func(t *testing.T) {
		cmd := NewShellCmd(context.Background(), "ping baidu.com -c 3")
		cmdResult, err := cmd.Run()
		assert.Nil(t, err)

		stdout, err := cmdResult.GetStdout()
		assert.Nil(t, err)
		lines := strings.Split(string(stdout), "\n")
		assert.Equal(t, 9, len(lines))

		stdout, err = cmdResult.GetStdout()
		if err != nil {
			panic(err)
		}
		fmt.Println(stdout)
	})

	t.Run("wait", func(t *testing.T) {
		cmd := NewShellCmd(context.Background(), "ping baidu.com -c 3")
		result, err := cmd.Run()
		assert.Nil(t, err)
		fmt.Println("run...")

		stdout, stderr, err := result.Wait()
		assert.Nil(t, err)
		fmt.Println(string(stdout))
		fmt.Println(string(stderr))
	})

	t.Run("scan", func(t *testing.T) {
		result, err := NewShellCmd(context.Background(), "sleep 3 && echo hello").Run()
		assert.Nil(t, err)

		err2 := result.ScanStdout(func(line string) {
			assert.Equal(t, "hello", line)
		})
		assert.Nil(t, err2)

		err3 := result.ScanStderr(func(line string) {
			assert.Equal(t, "", line)
		})
		assert.Nil(t, err3)
	})

	t.Run("double_read", func(t *testing.T) {
		result, err := NewShellCmd(context.Background(), "sleep 3 && echo hello").Run()
		assert.Nil(t, err)

		assert.Equal(t, "hello\n", result.GetStdoutStr())
		assert.Equal(t, "", result.GetStdoutStr())
	})
}

func TestGetStdout(t *testing.T) {
	cmdResult, err := NewShellCmd(context.Background(), "echo 'hello world'").Run()
	assert.Nil(t, err)

	assert.Equal(t, "hello world\n", cmdResult.GetStdoutStr())
}
