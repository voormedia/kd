# v2.1.0

* Add support for explicit default apps (`default: true`).

# v2.0.9

* Fix issue during initialization.

# v2.0.8

* Show more verbose output during pre/post-build steps.
* Warn if pre-build step contains reference to '.ssh' directory.

# v2.0.7

* Check if any SSH keys are exposed by `SSH_AUTH_SOCK`, and warn if there are none.

# v2.0.6

* Supply the app platform to docker when building with `kd build`.

# v2.0.4

* Add anti affinity to production template for `kd init`.

# v2.0.3

* Reduce default CPU reservation in production template for `kd init`.

# v2.0.2

* Do not attempt to forward SSH keys if `SSH_AUTH_SOCK` is unset.

# v2.0.1

* Fix `--tag` flag to work for `build` & `deploy` commands.

# v2.0.0

* Builds are executed directly with `docker buildx`, using BuildKit.
* Deploys are executed directly with `kubectl`.
* Configuration is resolved with the latest version of `kustomize`. This requires the configuration to be upgraded to version 2 with `kd upgrade`.
