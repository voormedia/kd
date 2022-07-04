package kubectl

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/voormedia/kd/pkg/config"
	"gopkg.in/yaml.v3"
)

func ApplyFromStdin(target *config.ResolvedTarget, input []byte) error {
	cmd := runCmdWithArgs(appendTargetArgs([]string{"apply", "-f", "-"}, target))

	cmd.Stdin = bytes.NewReader(input)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func RunForTarget(target *config.ResolvedTarget, args ...string) error {
	cmd := runCmdWithArgs(appendTargetArgs(args, target))

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func Version() (string, error) {
	bytes, err := capture("version", "--client", "--output", "json")
	if err != nil {
		return "", err
	}

	var version VersionDetails
	if err := yaml.Unmarshal(bytes, &version); err != nil {
		return "", err
	}

	return "v" + version.Major + "." + version.Minor, nil
}

type ClientVersion struct {
	Major string `yaml:"major,omitempty"`
	Minor string `yaml:"minor,omitempty"`
}

type VersionDetails struct {
	ClientVersion `yaml:"clientVersion,omitempty"`
}

func capture(args ...string) ([]byte, error) {
	cmd := runCmdWithArgs(args)

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

func runCmdWithArgs(args []string) *exec.Cmd {
	return exec.Command("kubectl", args...)
}

func appendTargetArgs(args []string, target *config.ResolvedTarget) []string {
	return append(args,
		"--context", target.Context,
		"--namespace", target.Namespace,
	)
}
