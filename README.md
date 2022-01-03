# Build and deploy apps to Kubernetes engine

## Prerequisites

1. A working Docker installation – https://store.docker.com/editions/community/docker-ce-desktop-mac
2. Google Cloud SDK `gcloud` – https://cloud.google.com/sdk/docs/quickstart-macos

## Installing

### Option 1 – prebuilt

1. Ensure you have `~/.bin` directory or similar that is in your `$PATH`
2. Install KD:

- When using macOS 11+ run: `curl -L $(curl -s https://api.github.com/repos/voormedia/kd/releases/latest | grep browser_download_url | grep big_sur_amd64 | cut -d '"' -f 4) -o ~/.bin/kd && chmod +x ~/.bin/kd`
- Otherwise run: `curl -L $(curl -s https://api.github.com/repos/voormedia/kd/releases/latest | grep browser_download_url | grep darwin_amd64 | cut -d '"' -f 4) -o ~/.bin/kd && chmod +x ~/.bin/kd`

3. Install Google Cloud credential helper:

- When using an Apple M1+ MacBook run: `curl -L https://github.com/GoogleCloudPlatform/docker-credential-gcr/releases/download/v2.1.0/docker-credential-gcr_darwin_arm64-2.1.0.tar.gz | tar -xzC ~/.bin docker-credential-gcr && chmod +x ~/.bin/docker-credential-gcr && docker-credential-gcr configure-docker`
- When using an Intel MacBook run: `curl -L https://github.com/GoogleCloudPlatform/docker-credential-gcr/releases/download/v2.1.0/docker-credential-gcr_darwin_amd64-2.1.0.tar.gz | tar -xzC ~/.bin docker-credential-gcr && chmod +x ~/.bin/docker-credential-gcr && docker-credential-gcr configure-docker`

### Option 2 – from source

1. Make sure you have a working `go` installation
2. Build KD from source: `go install github.com/voormedia/kd`
3. Install Google Cloud credential helper: `curl -L https://github.com/GoogleCloudPlatform/docker-credential-gcr/releases/download/v1.4.3/docker-credential-gcr_darwin_amd64-1.4.3.zip | funzip > ~/.bin/docker-credential-gcr && chmod +x ~/.bin/docker-credential-gcr && docker-credential-gcr configure-docker`

## Best practices for deploying

### Step 1 – adjust your app

- Configure the app to run with an application server (e.g. with `puma` for Rails).
- Make your app log to stdout/stderr instead of log files. Preferably in Google cloud compatible JSON.

### Step 2 – write a Dockerfile

- Make sure your image is small. Use a two-step build process. Use `.dockerignore` to exclude unused files.
- See https://github.com/voormedia/docker-base-images/tree/master/_examples for examples.

### Step 3 – configure deployment

- Run `kd init` and input the project details.
- Review the generated configuration and adjust where necessary.

### Step 4 – configure cloud services

- Create a PostgreSQL user and database if necessary. Use the same naming conventions as generated by `kd` in step 2. Create a secret from the service account key.

### Step 5 – deploy

- Use `kd deploy` to deploy to a target.
