# ocs-operator Test Harness for OpenShift Dedicated E2E Test Suite

This test harness is for validating [ocs-operator] deployed as an
add-on in an OpenShift Dedicated (OSD) cluster. _This code is NOT run
directly._

Compiled docker images of this test harness are pushed to [quay].
These are fetched by the [osde2e] tool and run inside the cluster.
Running the add-on using [osde2e] in the CI is covered in osde2e's
[add-on documentation]. However, that isn't helpful while developing
the add-on. This file contains those details and this repo supplies
other scripts and configuration files to enable running `osde2e`
against a cluster manually; optionally with this test harness.

Currently, the [ocs-operator] add-on needs to be installed manually.
The [deploy_ocs_on_osd.sh] script sets up the cluster environment to
enable the add-on to then be deployed either via the OSD dashboard or
using the script in [John Strunk's repo].

This repository contains:

- The [test harness] which is compiled into a docker image and run via [osde2e].
- [envrc] for configuring the environment variables for running [osde2e]. I normally use it in a custom `.envrc` file and load/reload it automatically using [direnv]. There's an example in [Custom .envrc for use with direnv].
- [deploy_ocs_on_osd.sh] script which can:
  - fetch cluster details into files that are loaded via [envrc]
  - prepare the cluster for the deployment of the add-on
- [auths.json.template] which is used with [deploy_ocs_on_osd.sh] to use custom auth tokens for registries.
- [osde2e_addons_config.yaml] which configures the `osde2e` tool to run the add-on test suite.


## deploy_ocs_on_osd script

***NOTE: The script assumes that only one cluster is active. It
literally uses only the first cluster in the list and ignores the
rest. It also reads and writes files in $PWD.***

### Prerequisites:

- `ocm` is in `$PATH` and is logged in.
- `auths.json` is populated using [auths.json.template] present in `$PWD` for custom pull secrets. Skip if not needed.

### Usage

The script can be used in two modes:

- _--details-only_: Fetch the cluster details and populate the files needed by [envrc]. This will enable [osde2e] to run.
- _--prepare_: Prepare the cluster for ocs-operator deployment. Always re-populates the [envrc] files.

In either of the modes, the script will write the following files to `$PWD`:
- _cluster_id_
- _cluster_name_
- _admin_: kubeadmin password
- _kubeconfig_: Point `$KUBECONFIG` its path to access the cluster with `oc`.
- _secret.json_: Pull secrets fetched from the cluster.
- _pull-secret-<cluster-id>.json_: Pull secrets compiled by merging the `auths.json` into `secret.json`; pushed to the cluster during `--prepare`.

### Preparing the cluster and deploying ocs-operator

- Create an OSD cluster with 9 worker nodes and 4 load balancers.
- Run the script:
  - Use `--details-only` if you wish simply to run [osde2e] against the cluster without deploying the add-on.
  - Use `--prepare` to also prepare the cluster to be able to deploy the ocs-operator.
- _If the script exits with an error due to desired cluster state not being reached in time, just re-run the script. It is idempotent and at the most only overwrite the cluster details stored in the files._
- Pull [John Strunk's repo] and deploy ocs-operator.


## Building the container image

### Pre-requisites

- Properly configured `go` workspace. While `$GOPATH` is mandatory, I also set `GO111MODULES=true` and `GOROOT=$(go env GOROOT)`. Check the example `.envrc` in the last section.
- Running `docker` daemon.
- `docker` client authenticated to the registry.
- Run `docker run hello-world` to check that `docker` is functional.

***Don't run `go tidy`. It will fetch the latest versions of some modules and break the suite.***

### Building the image

```bash
IMAGE_NAME="ocs-operator-test-harness"
IMAGE_TAG="0.01"
IMAGE_REPO="quay.io/mkarnikredhat/ocs-operator-osde2e-test-harness"
docker build -t "$IMAGE_NAME":"$IMAGE_TAG"
docker tag "$IMAGE_NAME":"$IMAGE_TAG" "$IMAGE_REPO":"$IMAGE_TAG"
docker push "$IMAGE_REPO":"$IMAGE_TAG"
```

## Running osde2e Manually

### Pre-requisites

- Properly configure the `go` workspace. While `$GOPATH` is mandatory, I also set `GO111MODULES=true` and `GOROOT=$(go env GOROOT)`.
- Checkout [osde2e] and compile the `odse2e` tool by running `make build`.
- Edit [envrc] to set the correct values for:
  - _OSD_PROJECT_DIR_: directory where all the cluster details files reside; in case it's not `$PWD`. **Use absolute path.**
  - _OCS_ADDON_TEST_HARNESS_: Container image repository for the test harness. Default should point to [quay].
  - _OCS_ADDON_TEST_HARNESS_TAG_: Image tag to pull. Currently there's no `latest` tag on the repo; but this will be updated to it as soon as there is and shouldn't need changing.
- The repo ships with [osde2e_addons_config.yaml] which will run only the addons suite by default. Update this to add any other suites to be run.

### Running the test suite

***Run the deploy script and re-load envrc each time the cluster is re-deployed.***

- Run `./deploy_ocs_on_osd.sh --details-only` to populate the cluster details files.
- `source envrc` to load the cluster details into environment variables and configure `osde2e`.
- `~/git/openshift/osde2e/out/osde2e test --custom-config osde2e_addons_config.yaml`. Obviously, use the correct paths.
- The output should be in `"$ARTIFACTS/install/junit-ocs-operator.xml"`.


## Custom .envrc for use with direnv

```
export GO111MODULE=on
export GOROOT="$(go env GOROOT)"
source envrc
```

I source another file which sets up `$KUBECONFIG` should I want to
run `oc`.


[ocs-operator]:https://github.com/openshift/ocs-operator
[osde2e]:https://github.com/openshift/osde2e
[John Strunk's repo]:https://github.com/JohnStrunk/ocp-rook-ceph
[quay]:https://quay.io/repository/mkarnikredhat/ocs-operator-osde2e-test-harness
[add-on documentation]:https://github.com/openshift/osde2e/blob/main/docs/Addons.md
[direnv]:https://direnv.net/
[deploy_ocs_on_osd.sh]:deploy_ocs_on_osd.sh
[test harness]:ocs_operator_test_harness_suite_test.go
[envrc]:envrc
[osde2e_addons_config.yaml]:osde2e_addons_config.yaml
