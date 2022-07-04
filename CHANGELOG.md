# v2.0.0

* Builds are executed directly with `docker buildx`, using BuildKit.
* Deploys are executed directly with `kubectl`.
* Configuration is resolved with the latest version of `kustomize`. This requires the configuration to be upgraded to version 2 with `kd upgrade`.
