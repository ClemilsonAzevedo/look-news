package feed

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

type req struct {
	Query    string    `json:"query"`
	Articles []Article `json:"articles"`
}

type response struct {
	Relevant []string `json:"relevant"`
}

type Filter struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	query  string
	mu     sync.Mutex
}

func Start(python, scriptPath, query string) (*Filter, error) {
	cmd := exec.Command(python, scriptPath)
	cmd.Dir = filepath.Dir(scriptPath)
	cmd.Stderr = os.Stderr // Colocar os logs no terminal!!

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Filter{
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewReader(stdout),
		query:  query,
	}, nil
}

func (f *Filter) ApplyFilter(arts []Article) ([]Article, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	body, err := json.Marshal(req{Query: f.query, Articles: arts})
	if err != nil {
		return nil, err
	}

	if _, err := f.stdin.Write(append(body, '\n')); err != nil {
		return nil, fmt.Errorf("stdin write in python: %w", err)
	}

	line, err := f.stdout.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("stdout read in python: %w", err)
	}

	var r response
	if err := json.Unmarshal(line, &r); err != nil {
		return nil, fmt.Errorf("invalid response: %w", err)
	}

	relevantSet := make(map[string]struct{}, len(r.Relevant))
	for _, link := range r.Relevant {
		relevantSet[link] = struct{}{}
	}

	var out []Article
	for _, article := range arts {
		if _, ok := relevantSet[article.Link]; ok {
			out = append(out, article)
		}
	}

	return out, nil
}

func (f *Filter) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.stdin != nil {
		_ = f.stdin.Close()
	}

	done := make(chan error, 1)
	go func() { done <- f.cmd.Wait() }()
	select {
	case err := <-done:
		return err
	case <-time.After(10 * time.Second):
		_ = f.cmd.Process.Kill()
		return fmt.Errorf("timed out waiting for process to finish")
	}
}
