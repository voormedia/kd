package kubectl

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/voormedia/kd/pkg/config"
	"gopkg.in/yaml.v3"
)

func ApplyFromStdin(target *config.ResolvedTarget, input []byte) error {
	err := RunForTarget(input, target, "apply", "-f", "-")

	if err != nil {
		return err
	}

	return nil
}

func RunForTarget(stdin []byte, target *config.ResolvedTarget, args ...string) error {
	args = append(args,
		"--context", target.Context,
		"--namespace", target.Namespace,
	)

	cmd := exec.Command("kubectl", args...)

	cmd.Stdin = bytes.NewReader(stdin)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func Version() (string, error) {
	bytes, err := capture(nil, "version", "--client", "--output", "json")
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

func capture(stdin []byte, args ...string) ([]byte, error) {
	cmd := exec.Command("kubectl", args...)
	buf := &bytes.Buffer{}
	cmd.Stdin = bytes.NewReader(stdin)
	cmd.Stderr = os.Stderr
	cmd.Stdout = buf

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
