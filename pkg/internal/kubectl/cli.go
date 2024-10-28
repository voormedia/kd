package kubectl

import (
	"encoding/json"

	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/util"
	networking "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/version"
)

func ApplyFromStdin(log *util.Logger, target *config.ResolvedTarget, input []byte) error {
	return util.RunWithInput(log, input,
		"kubectl",
		"--context", target.Context,
		"--namespace", target.Namespace,
		"apply", "-f", "-")
}

func GetGCEIngresses(log *util.Logger, target *config.ResolvedTarget) ([]*networking.Ingress, error) {
	bytes, err := util.Capture(log,
		"kubectl",
		"--context", target.Context,
		"--namespace", target.Namespace,
		"get", "ingress",
		"--output", "json")

	if err != nil {
		return nil, err
	}

	var list networking.IngressList
	if err := json.Unmarshal(bytes, &list); err != nil {
		return nil, err
	}

	var ingresses []*networking.Ingress
	for _, ingress := range list.Items {
		if ingress.Annotations["kubernetes.io/ingress.class"] == "gce" {
			ingresses = append(ingresses, &ingress)
		}
	}

	return ingresses, nil
}

func RunForTarget(log *util.Logger, target *config.ResolvedTarget, args ...string) error {
	args = append([]string{
		"--context", target.Context,
		"--namespace", target.Namespace,
	}, args...)

	return util.RunInteractively(log, "kubectl", args...)
}

func Version(log *util.Logger) (string, error) {
	bytes, err := util.Capture(log,
		"kubectl", "version",
		"--client",
		"--output", "json")

	if err != nil {
		return "", err
	}

	var version VersionDetails
	if err := json.Unmarshal(bytes, &version); err != nil {
		return "", err
	}

	return "v" + version.ClientVersion.Major + "." + version.ClientVersion.Minor, nil
}

type VersionDetails struct {
	ClientVersion version.Info `json:"clientVersion,omitempty"`
}
