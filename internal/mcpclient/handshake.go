package mcpclient

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

// DefaultHandshakeTimeout bounds how long `mcp status` waits for a configured MCP
// server to respond to a handshake before giving up.
const DefaultHandshakeTimeout = 5 * time.Second

type jsonrpcResponse struct {
	ID     json.RawMessage `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  *jsonrpcError   `json:"error"`
}

type jsonrpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// initializeRequest is a minimal MCP "initialize" JSON-RPC request, sufficient to
// verify that a server starts, authenticates, and speaks MCP over stdio.
func initializeRequest() []byte {
	return []byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"debricked-cli-status-check","version":"1.0.0"}}}`)
}

// CheckServer spawns command with args exactly as an MCP client would, sends a
// minimal initialize request over stdin, and waits up to timeout for a valid
// JSON-RPC response on stdout. This proves PATH resolution, the binary, and
// authentication all work end-to-end - not just that credentials look valid
// in the current process.
func CheckServer(ctx context.Context, command string, args []string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, args...) // #nosec G204 -- command/args come from the user's own MCP client config, not external input

	stdin, stdout, stderr, err := setupCommand(cmd)
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return handleStartError(command, err)
	}

	done := waitForCommand(cmd)
	response := readResponse(stdout)

	if err := sendHandshake(stdin, cmd, done); err != nil {
		return classifyFailure(stderr.String(), err)
	}

	result := waitForResponse(ctx, response)

	_ = cmd.Process.Kill()
	<-done

	if result != nil {
		return classifyFailure(stderr.String(), result)
	}

	return nil
}

func setupCommand(cmd *exec.Cmd) (io.WriteCloser, io.ReadCloser, *strings.Builder, error) {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to open stdin: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to open stdout: %w", err)
	}

	var stderr strings.Builder
	cmd.Stderr = &stderr

	return stdin, stdout, &stderr, nil
}

func handleStartError(command string, err error) error {
	if errors.Is(err, exec.ErrNotFound) {
		return fmt.Errorf("%q was not found on PATH", command)
	}

	return fmt.Errorf("failed to start %q: %w", command, err)
}

func waitForCommand(cmd *exec.Cmd) <-chan error {
	done := make(chan error, 1)

	go func() {
		done <- cmd.Wait()
	}()

	return done
}

func sendHandshake(
	stdin io.WriteCloser,
	cmd *exec.Cmd,
	done <-chan error,
) error {
	_, err := stdin.Write(append(initializeRequest(), '\n'))
	if err == nil {
		return nil
	}

	_ = cmd.Process.Kill()
	<-done

	return fmt.Errorf("failed to send handshake request: %w", err)
}

func readResponse(stdout io.Reader) <-chan error {
	response := make(chan error, 1)

	go func() {
		scanner := bufio.NewScanner(stdout)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}

			var resp jsonrpcResponse
			if err := json.Unmarshal([]byte(line), &resp); err != nil {
				response <- fmt.Errorf("received a malformed response: %w", err)

				return
			}

			if resp.Error != nil {
				response <- fmt.Errorf(
					"server returned an error: %s",
					resp.Error.Message,
				)

				return
			}

			response <- nil

			return
		}

		if err := scanner.Err(); err != nil {
			response <- fmt.Errorf("failed to read server output: %w", err)

			return
		}

		response <- errors.New("server closed its output before responding")
	}()

	return response
}

func waitForResponse(ctx context.Context, response <-chan error) error {
	select {
	case <-ctx.Done():
		return errors.New("timed out waiting for a response")
	case result := <-response:
		return result
	}
}

// classifyFailure turns a raw failure plus captured stderr into an actionable message,
// recognizing the specific errors `debricked mcp start` (see start.go) prints.
func classifyFailure(stderrOutput string, cause error) error {
	stderrOutput = strings.TrimSpace(stderrOutput)

	switch {
	case strings.Contains(stderrOutput, "no access token found"):
		return errors.New("no access token found; run `debricked auth login` or pass --access-token")

	case strings.Contains(stderrOutput, "access token rejected"):
		return fmt.Errorf("access token rejected: %s", stderrOutput)

	case stderrOutput != "":
		return fmt.Errorf("%w (stderr: %s)", cause, stderrOutput)

	default:
		return cause
	}
}
