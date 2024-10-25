package kustomize

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/util"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/provider"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

func GetResources(log *util.Logger, app *config.ResolvedApp, target *config.ResolvedTarget, digest string) ([]byte, error) {
	fSys := filesys.MakeFsOnDisk()

	kust := krusty.MakeKustomizer(
		krusty.MakeDefaultOptions(),
	)

	log.Debug("Running kustomize build for", filepath.Join(app.Path, target.Path))
	res, err := kust.Run(fSys, filepath.Join(app.Path, target.Path))
	if err != nil {
		return nil, err
	}

	resmapFactorty := resmap.NewFactory(provider.NewDepProvider().GetResourceFactory())
	all, err := resmapFactorty.NewResMapFromBytes(namespace(target))
	if err != nil {
		return nil, err
	}

	err = all.AbsorbAll(res)
	if err != nil {
		return nil, err
	}

	yml, err := all.AsYaml()
	if err != nil {
		return nil, err
	}

	// Search and replace in the generated YAML to set the actual deployment
	// image. This is the reason we cannot use 'kubectl -k ..' directly.
	url := app.RepositoryWithDigest(digest)
	buf := bytes.NewBuffer(bytes.Replace(yml, []byte(" image: "+app.Name), []byte(" image: "+url), -1))

	return buf.Bytes(), nil
}

func namespace(target *config.ResolvedTarget) []byte {
	var out bytes.Buffer
	namespaceTmpl.Execute(&out, target)
	return out.Bytes()
}

var namespaceTmpl = template.Must(template.New("namespace.yaml").Parse(
	`apiVersion: v1
kind: Namespace
metadata:
  name: {{.Namespace}}
`))
