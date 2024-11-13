# v2.5.4

* Combine build and push in a single step.

# v2.5.3

* Increase compatibility with older docker clients.

# v2.5.2

* Display error output if pushing to remote registry fails.

# v2.5.1

* Fix `kd build` to include output to an image even if this is not default.

# v2.5.0

* Automatically use remote cache for builds if supported by the builder (requires containerd).
* Do not automatically flush CDN cache on deploy, but only when explicitly requested with `--clear-cdn-cache`.

# v2.4.0

* Add support for `skipBuild` configuration option for apps, which allows deploying apps without building an app-specific image.

# v2.3.0

* Introduce aliases `kbuild` for `kd build`, `kdeploy` for `kd deploy` and `kctl` for `kd ctl`.

# v2.2.1

* Fix issues that prevented the `kd ctl` command from working.

# v2.2.0

* Automatically flush any CDN cache on deployment.
* Add verbose flag.
* Update internal dependencies.

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
