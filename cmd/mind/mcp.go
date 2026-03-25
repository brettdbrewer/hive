package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// mcpConfig is the JSON config for an MCP server subprocess.
type mcpConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
}

// mcpTool mirrors the tool definition returned by the MCP server.
type mcpTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema any    `json:"inputSchema"`
}

// mcpToolResult mirrors the tool result returned by the MCP server.
type mcpToolResult struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	IsError bool `json:"isError,omitempty"`
}

// mcpClient manages a subprocess running an MCP server.
type mcpClient struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	reader *bufio.Scanner
	tools  []mcpTool
	nextID int64
}

// startMCPClient starts the MCP server subprocess, performs the handshake,
// and returns a client ready to call tools.
func startMCPClient(ctx context.Context, cfg mcpConfig) (*mcpClient, error) {
	cmd := exec.CommandContext(ctx, cfg.Command, cfg.Args...)

	// Build env: inherit current env, then overlay cfg.Env.
	cmd.Env = os.Environ()
	for k, v := range cfg.Env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	cmd.Stderr = os.Stderr // forward server stderr to our stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start: %w", err)
	}

	mc := &mcpClient{
		cmd:    cmd,
		stdin:  stdin,
		reader: bufio.NewScanner(stdout),
		nextID: 1,
	}
	mc.reader.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	// MCP handshake: initialize.
	if err := mc.initialize(); err != nil {
		cmd.Process.Kill()
		return nil, fmt.Errorf("initialize: %w", err)
	}

	// List tools.
	tools, err := mc.listTools()
	if err != nil {
		cmd.Process.Kill()
		return nil, fmt.Errorf("tools/list: %w", err)
	}
	mc.tools = tools

	return mc, nil
}

func (mc *mcpClient) close() {
	mc.stdin.Close()
	mc.cmd.Wait()
}

func (mc *mcpClient) sendRequest(method string, params any) (json.RawMessage, error) {
	id := mc.nextID
	mc.nextID++

	req := map[string]any{
		"jsonrpc": "2.0",
		"id":      id,
		"method":  method,
	}
	if params != nil {
		req["params"] = params
	}

	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	b = append(b, '\n')

	if _, err := mc.stdin.Write(b); err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}

	// Read response lines until we get one with our ID.
	for mc.reader.Scan() {
		line := mc.reader.Bytes()
		var resp struct {
			JSONRPC string          `json:"jsonrpc"`
			ID      json.RawMessage `json:"id"`
			Result  json.RawMessage `json:"result"`
			Error   *struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal(line, &resp); err != nil {
			continue // skip non-JSON lines
		}
		var respID int64
		if err := json.Unmarshal(resp.ID, &respID); err != nil || respID != id {
			continue // not our response
		}
		if resp.Error != nil {
			return nil, fmt.Errorf("rpc error %d: %s", resp.Error.Code, resp.Error.Message)
		}
		return resp.Result, nil
	}
	return nil, fmt.Errorf("connection closed before response for id=%d", id)
}

func (mc *mcpClient) sendNotification(method string, params any) error {
	req := map[string]any{
		"jsonrpc": "2.0",
		"method":  method,
	}
	if params != nil {
		req["params"] = params
	}
	b, _ := json.Marshal(req)
	b = append(b, '\n')
	_, err := mc.stdin.Write(b)
	return err
}

func (mc *mcpClient) initialize() error {
	params := map[string]any{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]any{},
		"clientInfo":      map[string]any{"name": "mind", "version": "1.0.0"},
	}
	_, err := mc.sendRequest("initialize", params)
	if err != nil {
		return err
	}
	// Send initialized notification (no response expected).
	return mc.sendNotification("notifications/initialized", nil)
}

func (mc *mcpClient) listTools() ([]mcpTool, error) {
	result, err := mc.sendRequest("tools/list", nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Tools []mcpTool `json:"tools"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, err
	}
	return resp.Tools, nil
}

func (mc *mcpClient) callTool(ctx context.Context, name string, args map[string]any) (mcpToolResult, error) {
	params := map[string]any{
		"name":      name,
		"arguments": args,
	}
	result, err := mc.sendRequest("tools/call", params)
	if err != nil {
		return mcpToolResult{}, err
	}
	var toolResult mcpToolResult
	if err := json.Unmarshal(result, &toolResult); err != nil {
		return mcpToolResult{}, err
	}
	return toolResult, nil
}
