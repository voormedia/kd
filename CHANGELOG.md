# v2.0.2

* Do not attempt to forward SSH keys if `SSH_AUTH_SOCK` is unset.

# v2.0.1

* Fix `--tag` flag to work for `build` & `deploy` commands.

# v2.0.0

* Builds are executed directly with `docker buildx`, using BuildKit.
* Deploys are executed directly with `kubectl`.
* Configuration is resolved with the latest version of `kustomize`. This requires the configuration to be upgraded to version 2 with `kd upgrade`.
