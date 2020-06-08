package cmd

import (
	"os/exec"
	"strings"
	"testing"
)

// There is not any mocking. Do not run tests when use the this tool.

func TestSetCommand_WhenRunsWithEmptyArgs_ReturnsOutput(t *testing.T) {
	cmd := exec.Command("go", "run", "../main.go", "set",
		"--output", "")

	out, err := cmd.Output()
	if err != nil {
		t.Errorf(err.Error())
	}

	want := "output parameter error please check and re run"
	got := string(out)
	if !strings.Contains(got, want) {
		t.Errorf("Unexpected data.\nGot: %s\nExpected: %s", got, want)
	}
}

func TestSetCommand_WhenRunsWithCorrectArgs_ReturnsOutput(t *testing.T) {
	cmd := exec.Command("go", "run", "../main.go", "set",
		"--output", "/Users/furkan.bozdag/Desktop")

	out, err := cmd.Output()
	if err != nil {
		t.Errorf(err.Error())
	}

	want := ""
	got := string(out)
	if !strings.Contains(got, want) {
		t.Errorf("Unexpected data.\nGot: %s\nExpected: %s", got, want)
	}
}