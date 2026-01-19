package e2e

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func GetCommand(args ...string) (*exec.Cmd, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, binName)

	cmd := exec.Command(path, args...)
	return cmd, nil
}

func TestStdin(t *testing.T) {
	cmd, err := GetCommand()
	if err != nil {
		t.Fatalf("failed to get command: %v", err)
	}

	output := &bytes.Buffer{}

	cmd.Stdin = strings.NewReader("one two three\n")
	cmd.Stdout = output

	if err := cmd.Run(); err != nil {
		t.Fatal("Failed to run binary")
	}

	expected := " 1 3 14\n"

	if expected != output.String() {
		t.Errorf("unexpected output: got %q, want %q", output.String(), expected)
		t.Fail()
	}
}

func TestSingleFile(t *testing.T) {
	file, err := os.CreateTemp("", "counter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	_, err = file.WriteString("foo bar baz\none two three")
	if err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	err = file.Close()
	if err != nil {
		t.Fatalf("failed to close file: %v", err)
	}

	cmd, err := GetCommand(file.Name())
	if err != nil {
		t.Fatalf("failed to get command: %v", err)
	}
	output := &bytes.Buffer{}
	cmd.Stdout = output

	if err := cmd.Run(); err != nil {
		t.Fatal("Failed to run binary")
	}

	expected := fmt.Sprintf(" 1 6 25 %s\n 1 6 25 total\n", file.Name())

	if expected != output.String() {
		t.Errorf("unexpected output: got %q, want %q", output.String(), expected)
		t.Fail()
	}
}

func TestNoExist(t *testing.T) {
	cmd, err := GetCommand("noexist.txt")
	if err != nil {
		t.Fatalf("failed to get command: %v", err)
	}
	output := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = output
	cmd.Stderr = stderr

	expectedOut := "open noexist.txt: no such file or directory\n 0 0 0 total\n"
	expectedErr := ""

	err = cmd.Run()
	if err == nil {
		t.Errorf("expected error, got nil")
		t.Fail()
	}

	if expectedOut != output.String() {
		t.Errorf("unexpected output: got %q, want %q", output.String(), expectedOut)
		t.Fail()
	}

	if expectedErr != stderr.String() {
		t.Errorf("unexpected error output: got %q, want %q", stderr.String(), expectedErr)
		t.Fail()
	}
}
