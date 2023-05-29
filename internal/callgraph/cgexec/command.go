package cgexec

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ICommand interface {
	CombinedOutput() ([]byte, error)
	Start() error
	Wait() error
	GetProcess() *os.Process
	SetStdErr()
	SetStdout()
	GetArgs() []string
	GetDir() string
}

type Command struct {
	osCmd *exec.Cmd
}

func (cmd Command) SetStderr(stderr *bytes.Buffer) {
	cmd.osCmd.Stderr = stderr
}

func (cmd Command) SetStdout(stdout *bytes.Buffer) {
	cmd.osCmd.Stdout = stdout
}

func (cmd Command) GetArgs() []string {
	return cmd.osCmd.Args
}

func (cmd Command) CombinedOutput() ([]byte, error) {
	return cmd.osCmd.CombinedOutput()
}

func (cmd Command) Start() error {
	return cmd.osCmd.Start()
}

func (cmd Command) Wait() error {
	return cmd.osCmd.Wait()
}

func (cmd Command) GetProcess() *os.Process {
	return cmd.osCmd.Process
}

func (cmd Command) GetDir() string {
	return cmd.osCmd.Dir
}

func RunCommand(osCmd *exec.Cmd, ctx IContext) error {
	cmd := Command{osCmd}
	args := strings.Join(cmd.GetArgs(), " ")
	var stdoutBuf, stderrBuf bytes.Buffer
	// cmd.Stderr = &stderrBuf
	cmd.SetStderr(&stderrBuf)
	var err error
	var outputCmd []byte
	if ctx == nil {
		outputCmd, err = cmd.CombinedOutput()
		if _, ok := err.(*exec.ExitError); ok {
			err = fmt.Errorf("Command '%s' executed in folder '%s' gave the following error:\n%s", args, cmd.GetDir(), outputCmd)
		}

		return err
	}

	// cmd.Stdout = &stdoutBuf
	cmd.SetStdout(&stdoutBuf)

	// Start the external process
	if err := cmd.Start(); err != nil {

		return err
	}

	// Channel to signal when the external process completes
	done := make(chan error, 1)

	// Goroutine to wait for the process to complete
	go func() {
		err := cmd.Wait()
		if err != nil {
			err = fmt.Errorf("Command '%s' executed in folder '%s' gave the following error: \n%s\n%s", args, cmd.GetDir(), stdoutBuf.String(), stderrBuf.String())
		}

		done <- err
	}()

	select {
	case <-ctx.Context().Done():

		if ctx.Context().Err() == context.DeadlineExceeded {
			// Context timeout occurred before the process started
			return fmt.Errorf("Timeout error: Set timeout duration for Callgraph jobs reached")
		}

		// The context was canceled, handle cancellation if needed
		// Send a signal to the process to terminate
		process := cmd.GetProcess()
		if process != nil {
			err := process.Signal(os.Interrupt)
			if err != nil {
				return err
			}
		}

		// Wait for the process to exit
		<-done

		return fmt.Errorf("Timeout error: Set timeout duration for Callgraph jobs reached")
	case err := <-done:
		// The external process completed before the context was canceled
		return err
	}
}

func MakeCommand(workingDir string, path string, args []string, ctx IContext) *exec.Cmd {
	var cmd *exec.Cmd

	if ctx == nil {
		cmd = &exec.Cmd{
			Path: path,
			Args: args,
			Dir:  workingDir,
		}
	} else {
		command := args[0]
		arguments := args[1:]
		cmd = exec.CommandContext(ctx.Context(), command, arguments...)
		cmd.Path = path
		cmd.Dir = workingDir
	}

	return cmd
}
