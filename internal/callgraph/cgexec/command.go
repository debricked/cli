package cgexec

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunCommand(cmd *exec.Cmd, ctx IContext) error {
	if ctx == nil {
		_, err := cmd.Output()

		return err
	}

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
			args := strings.Join(cmd.Args, " ")
			err = fmt.Errorf("Command '%s' executed in folder '%s' gave the following error: %s", args, cmd.Dir, err.Error())
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
		if cmd.Process != nil {
			err := cmd.Process.Signal(os.Interrupt)
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
