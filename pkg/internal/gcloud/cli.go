package gcloud

import (
	"bytes"
	"encoding/json"

	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/util"
)

func ScheduleCDNCacheFlush(log *util.Logger, target *config.ResolvedTarget, annotations map[string]string) (bool, error) {
	backendData := annotations["ingress.kubernetes.io/backends"]
	urlMap := annotations["ingress.kubernetes.io/url-map"]

	if backendData == "" || urlMap == "" {
		return false, nil
	}

	var backends map[string]string
	json.Unmarshal([]byte(backendData), &backends)

	hasCdn := false
	for backend, _ := range backends {
		output, err := util.Capture(log,
			"gcloud",
			"compute", "backend-services", "describe", backend,
			"--global",
			"--project", target.GCPProject())

		if err != nil {
			return false, err
		}

		if bytes.Contains(output, []byte("enableCDN: true")) {
			hasCdn = true
			break
		}
	}

	if !hasCdn {
		return false, nil
	}

	err := util.Run(log,
		"gcloud",
		"compute", "url-maps", "invalidate-cdn-cache", urlMap,
		"--global",
		"--path", "/*",
		"--async",
		"--quiet",
		"--no-user-output-enabled",
		"--project", target.GCPProject())

	if err != nil {
		return false, err
	}

	return true, nil
}
