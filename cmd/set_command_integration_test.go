package cmd

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestSetCommand_WhenRunsWithEmptyArgs_ReturnsOutput(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

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
	if testing.Short() {
		t.Skip()
	}

	userHomeDir, err := os.UserHomeDir()

	if err != nil {
		t.Error(err)
	}

	cmd := exec.Command("go", "run", "../main.go", "set",
		"--output", userHomeDir)

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