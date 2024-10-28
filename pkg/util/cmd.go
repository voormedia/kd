package util

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

var Verbose bool

func Run(log *Logger, name string, args ...string) error {
	log.Debug("Executing:", name, strings.Join(args, " "))

	cmd := exec.Command(name, args...)
	cmd.Stdin = bytes.NewReader([]byte{})
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func RunInteractively(log *Logger, name string, args ...string) error {
	log.Debug("Executing interactively:", name, strings.Join(args, " "))

	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func RunWithInput(log *Logger, input []byte, name string, args ...string) error {
	log.Debug("Executing with input:", name, strings.Join(args, " "))

	cmd := exec.Command(name, args...)
	cmd.Stdin = bytes.NewReader(input)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func Capture(log *Logger, name string, args ...string) ([]byte, error) {
	log.Debug("Executing and capturing output:", name, strings.Join(args, " "))

	cmd := exec.Command(name, args...)
	buf := &bytes.Buffer{}
	cmd.Stdin = bytes.NewReader([]byte{})
	cmd.Stderr = os.Stderr
	cmd.Stdout = buf

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
