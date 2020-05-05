package exec

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

// RunCommand - run local command
func RunCommand(name string, args ...string) (string, error) {
	var out, berr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &out
	cmd.Stderr = &berr

	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("err:%v berr:%v", err, berr.String())
	}

	return out.String() + berr.String(), nil

}

// RunCommandContext - run context Command
func RunCommandContext(ctx context.Context, name string, args ...string) (string, error) {
	var out, berr bytes.Buffer
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = &out
	cmd.Stderr = &berr

	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("err:%v berr:%v", err, berr.String())
	}

	return out.String() + berr.String(), nil
}

// RunShellCommand - run shell command
func RunShellCommand(shell string) (string, error) {
	var out, berr bytes.Buffer
	cmd := exec.Command("/bin/sh", "-c", shell)
	cmd.Stdout = &out
	cmd.Stderr = &berr

	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("err:%v berr:%v", err, berr.String())
	}

	return out.String() + berr.String(), nil
}
