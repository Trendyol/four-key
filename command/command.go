package command

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

const DefaultFourKeyDirName = "four-key"

func (c *Commander) Info(args ...interface{}) string {
	f := c.color("\033[1;36m%s\033[0m")
	return f(args)
}

func (c *Commander) Good(args ...interface{}) string {
	f := c.color("\033[1;32m%s\033[0m")
	return f(args)

}

func (c *Commander) Fatal(args ...interface{}) string {
	f := c.color("\033[1;31m%s\033[0m")
	return f(args)

}

func (c *Commander) Warn(args ...interface{}) string {
	f := c.color("\033[1;33m%s\033[0m")
	return f(args)

}

var c Commander

func ACommander() *Commander {
	return &c
}

type Commander struct {
}

type ICommand interface {
	Command(cmd string, p string) error
	GetFourKeyPath() string
	GetRepositoriesPath(cloneDir string) string
	GetUserHomeDir() string
	Info(...interface{}) string
	Warn(...interface{}) string
	Fatal(...interface{}) string
	Good(...interface{}) string
	Open(path string) error
}

func (c *Commander) Command(command string, p string) error {
	// Prepare the command to execute.
	cmd := exec.Command("sh", "-c", strings.TrimSuffix(command, "\n"))

	// Set the correct output device.
	cmd.Dir = p

	// Execute the command and return the error.
	output, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(output))
		return err
	}

	return nil
}

func (c *Commander) GetFourKeyPath() string {
	r := c.GetUserHomeDir()

	p := path.Join(r, DefaultFourKeyDirName)
	err := os.Chdir(p)

	if err != nil {
		log.Fatal(c.Fatal("four-key path not found! Error: %v", err))
	}

	return p
}

func (c *Commander) GetUserHomeDir() string {
	r, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err)
	}

	return r
}

func (c *Commander) GetRepositoriesPath(cloneDir string) string {
	r := c.GetUserHomeDir()

	p := path.Join(r, DefaultFourKeyDirName, cloneDir)
	err := os.Mkdir(p, os.ModePerm)

	log.Println(err)

	return p
}

func (c *Commander) color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func (c *Commander) Open(path string) error {
	if runtime.GOOS == "windows" {
		return c.Command("start .", path)
	} else {
		return c.Command("open .", path)
	}
}
