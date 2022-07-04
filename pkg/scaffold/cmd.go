package scaffold

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/util"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func Run(log *util.Logger) error {
	fs := &afero.Afero{Fs: afero.NewOsFs()}

	details, err := requestDetails(fs, log, os.Stdin, os.Stdout)
	if err != nil {
		return err
	}

	err = writeConfig(fs, details)
	if err != nil {
		return err
	}

	if len(details.Apps) == 0 {
		log.Success("Created example configuration without any apps")
	} else {
		log.Success("Created example configuration for", strings.Join(details.appNames(), ", "))
	}

	log.Note("Next:  1. Review and adjust configuration")
	log.Note("       2. Make sure your apps log to stdout/stderr")
	log.Note("       3. Let your apps respond to health checks at /healthz")
	log.Note("       4. Use kd to build and deploy")

	return nil
}

type app struct {
	Name string
	Path string
}

type appEnv struct {
	app
	Namespace   string
	Environment string
	Replicas    int
}

type details struct {
	ApiVersion uint
	Customer   string
	Project    string
	Context    string
	Apps       []app
}

func (d details) appNames() (names []string) {
	names = make([]string, len(d.Apps))
	for i, app := range d.Apps {
		names[i] = app.Name
	}
	return
}

func requestDetails(afs *afero.Afero, log *util.Logger, in io.Reader, out io.Writer) (*details, error) {
	projects, err := findProjects()
	if err != nil {
		return nil, err
	}

	contexts, err := findContexts()
	if err != nil {
		return nil, err
	}

	apps, err := findApps(afs)
	if err != nil {
		return nil, err
	}

	data := &details{
		ApiVersion: config.LatestVersion,
		Apps:       apps,
	}

	if len(apps) == 0 {
		log.Warn("Could not find any apps; are you missing a Dockerfile?")
		log.Warn("Create a Dockerfile and rerun init to configure your app")
	}

	log.Note("Enter a few project details")

	qs := []*survey.Question{
		{
			Name:      "customer",
			Validate:  survey.Required,
			Prompt:    &survey.Input{Message: "Namespace (e.g. customer name):"},
			Transform: survey.TransformString(util.Slugify),
		},
		{
			Name:     "project",
			Validate: survey.Required,
			Prompt: &survey.Select{
				Message: "Select Google cloud project id:",
				Options: projects,
			},
		},
		{
			Name:     "context",
			Validate: survey.Required,
			Prompt: &survey.Select{
				Message: "Select Kubernetes cluster context:",
				Options: contexts,
			},
		},
	}

	err = survey.Ask(qs, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func findProjects() ([]string, error) {
	cmd := exec.Command("gcloud", "projects", "list", "--uri")
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		return nil, errors.Errorf("Failed to get projects: %s", errOut.String())
	}

	projects := strings.Split(out.String(), "\n")
	for i, proj := range projects[:len(projects)-1] {
		projects[i] = filepath.Base(proj)
	}

	return projects, nil
}

func findContexts() ([]string, error) {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	loader := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
	loadedConfig, err := loader.Load()
	if err != nil {
		return nil, err
	}

	contexts := make([]string, len(loadedConfig.Contexts))
	i := 0
	for ctx := range loadedConfig.Contexts {
		contexts[i] = ctx
		i++
	}

	return contexts, nil
}

func findApps(afs *afero.Afero) ([]app, error) {
	root, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return nil, err
	}

	var apps []app
	afs.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.Name() == "Dockerfile" {
			path := filepath.Dir(path)
			name := filepath.Base(path)

			path = strings.Replace(path, root+"", "", 1)
			path = strings.Replace(path, "/", "", 1)
			if path == "" {
				path = "."
			}

			apps = append(apps, app{
				Path: path,
				Name: name,
			})
		}

		return nil
	})

	return apps, nil
}

var kdeploy = template.Must(template.New(config.ConfigName).Parse(
	`# Check version compatibility with kd
version: {{.ApiVersion}}

# Private docker registry to push images to
registry: eu.gcr.io/{{.Project}}/{{.Customer}}

# List of apps to build
apps:
{{- range .Apps}}
- name: {{.Name}}
  path: {{.Path}}
{{- end}}

# List of available deployment targets
targets:
- name: acceptance
  alias: acc
  context: {{.Context}}
  namespace: {{.Customer}}-acc
  path: config/deploy/acceptance

- name: production
  alias: prd
  context: {{.Context}}
  namespace: {{.Customer}}-prd
  path: config/deploy/production
`))

var bseManifest = template.Must(template.New("kustomization.yaml").Parse(
	`# List of base resources
resources:
- deployment.yaml
- service.yaml
- ingress.yaml
`))

var envManifest = template.Must(template.New("kustomization.yaml").Parse(
	`# List of patches to apply (in order) for this environment
patches:
- deployment.yaml
- ingress.yaml

# Patches are applied to base resources
resources:
- ../_base
`))

var bseService = template.Must(template.New("service.yaml").Parse(
	`# Service defines the common entrypoint for multiple application pods
apiVersion: v1
kind: Service
metadata:
  name: {{.Name}}
  labels:
    app: {{.Name}}

spec:
  type: NodePort
  ports:
  - port: 80
    name: http
  selector:
    app: {{.Name}}
`))

var bseIngress = template.Must(template.New("ingress.yaml").Parse(
	`# Defines a cloud load balancer for HTTP + HTTPS
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ingress
  annotations:
    # Causes cert-manager to provision a TLS certificate with Let's encrypt.
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    acme.cert-manager.io/http01-edit-in-place: "true"
    kubernetes.io/tls-acme: "true"

    # This must be set in order for kube-lego to create proper forwarding
    # rules for the Let's Encrypt challenge/response endpoints.
    kubernetes.io/ingress.class: "gce"
    kubernetes.io/ingress.allow-http: "true"
    kubernetes.io/ingress.force-ssl-redirect: "true"

spec:
  defaultBackend:
    service:
      name: {{.Name}}
      port:
        number: 80
`))

var envIngress = template.Must(template.New("ingress.yaml").Parse(
	`# Defines a cloud load balancer for HTTP + HTTPS
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ingress
  annotations:
    kubernetes.io/ingress.global-static-ip-name: {{.Namespace}}

spec:
  tls:
  - secretName: {{.Name}}-tls
    hosts:
    # Replace this with a list of actual {{.Environment}} hostnames for this
    # application. Only primary hostnames should be added. Secondary hostnames
    # can be set up to redirect to a primary hostname, but do not need HTTPS.
    - {{if eq .Environment "acceptance"}}acceptance.{{end}}{{.Name}}.voormedia.com
`))

var bseDeployment = template.Must(template.New("deployment.yaml").Parse(
	`# Defines an app consisting of one or more identical pods
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Name}}

spec:
  selector:
    matchLabels:
      app: {{.Name}}

  revisionHistoryLimit: 5
  minReadySeconds: 5

  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0

  template:
    metadata:
      labels:
        app: {{.Name}}

    spec:
      containers:
      - name: {{.Name}}

        # This will be extrapolated on deploy by kd to refer to the
        # latest image that was pushed to the repository.
        image: {{.Name}}

        ports:
        - containerPort: 80
          name: http

        # Define environment variables that apply to all environments.
        env:
        - name: PORT
          value: "80"
        - name: RAILS_MAX_THREADS
          value: "16"

        readinessProbe:
          # The readiness probe is checked by Kubernetes, but is also adopted
          # as a cloud load balancer health check. Only HTTP GET is supported!
          httpGet:
            path: /healthz
            port: 80
          initialDelaySeconds: 2
          periodSeconds: 5
          timeoutSeconds: 1

      # PostgreSQL proxy container.
      # - name: cloudsql-proxy
      #   image: gcr.io/cloudsql-docker/gce-proxy:1.11
      #   command: [
      #     "/cloud_sql_proxy",
      #     "-instances=<INSTANCE-NAME>=tcp:5432",
      #     "-credential_file=/secrets/cloudsql/credentials.json"
      #   ]
      #
      #   volumeMounts:
      #   - name: cloudsql-instance-credentials
      #     mountPath: /secrets/cloudsql
      #     readOnly: true

      # To create this secret, run:
      #   kd kubectl <ENVIRONMENT> create secret generic
      #   cloudsql-instance-credentials --from-file=credentials.json=<KEY-FILE>
      # volumes:
      # - name: cloudsql-instance-credentials
      #   secret:
      #     secretName: cloudsql-instance-credentials
`))

var envDeployment = template.Must(template.New("deployment.yaml").Parse(
	`# Defines an app consisting of one or more identical pods
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Name}}

spec:
  replicas: {{.Replicas}}

  template:
    spec:
      containers:
      - name: {{.Name}}

        # Define environment variables specific to {{.Environment}}.
        env:
        - name: RACK_ENV
          value: {{.Environment}}

{{- if eq .Environment "production"}}
        resources:
          # Reserve resources for this application.
          requests:
            cpu: 100m
            memory: 250Mi

          # Limit CPU and memory usage of this application. If the pod uses
          # more than the specified amount of memory, it will be terminated.
          limits:
            cpu: 500m
            memory: 500Mi

{{else}}
        resources:
          # No minimum resources means this application will be scheduled in
          # {{.Environment}} on a best-effort basis.
          requests:
            cpu: 0m
            memory: 0Mi

          # Limit CPU and memory usage of this application. If the pod uses
          # more than the specified amount of memory, it will be terminated.
          limits:
            cpu: 500m
            memory: 500Mi
{{- end}}
`))

func writeConfig(afs *afero.Afero, details *details) error {
	fs := &fsMonad{afs: afs}
	fs.writeTemplate(config.ConfigName, kdeploy, details)

	for _, bseApp := range details.Apps {
		bsePath := filepath.Join(bseApp.Path, "config", "deploy", "_base")
		accPath := filepath.Join(bseApp.Path, "config", "deploy", "acceptance")
		prdPath := filepath.Join(bseApp.Path, "config", "deploy", "production")

		accApp := &appEnv{app: bseApp,
			Environment: "acceptance",
			Namespace:   details.Customer + "-acc",
			Replicas:    1,
		}

		prdApp := &appEnv{app: bseApp,
			Environment: "production",
			Namespace:   details.Customer + "-prd",
			Replicas:    2,
		}

		fs.mkdir(accPath)
		fs.mkdir(prdPath)

		fs.writeTemplate(filepath.Join(bsePath, "kustomization.yaml"), bseManifest, bseApp)
		fs.writeTemplate(filepath.Join(accPath, "kustomization.yaml"), envManifest, accApp)
		fs.writeTemplate(filepath.Join(prdPath, "kustomization.yaml"), envManifest, prdApp)

		fs.writeTemplate(filepath.Join(bsePath, "service.yaml"), bseService, bseApp)

		fs.writeTemplate(filepath.Join(bsePath, "ingress.yaml"), bseIngress, bseApp)
		fs.writeTemplate(filepath.Join(accPath, "ingress.yaml"), envIngress, accApp)
		fs.writeTemplate(filepath.Join(prdPath, "ingress.yaml"), envIngress, prdApp)

		fs.writeTemplate(filepath.Join(bsePath, "deployment.yaml"), bseDeployment, bseApp)
		fs.writeTemplate(filepath.Join(accPath, "deployment.yaml"), envDeployment, accApp)
		fs.writeTemplate(filepath.Join(prdPath, "deployment.yaml"), envDeployment, prdApp)
	}

	return fs.err
}

type fsMonad struct {
	err error
	afs *afero.Afero
}

func (fs *fsMonad) mkdir(path string) {
	if fs.err != nil {
		return
	}

	fs.err = fs.afs.MkdirAll(path, 0755)
	if fs.err != nil {
		return
	}
}

func (fs *fsMonad) writeTemplate(path string, tmpl *template.Template, data interface{}) {
	if fs.err != nil {
		return
	}

	var buf bytes.Buffer
	fs.err = tmpl.Execute(&buf, data)
	if fs.err != nil {
		return
	}

	fs.err = fs.afs.WriteFile(path, buf.Bytes(), 0644)
	if fs.err != nil {
		return
	}
}
