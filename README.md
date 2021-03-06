# chart-verifier

`chart-verifier` is a tool that certifies a Helm chart against a configurable list of checks; those checks can be
whitelisted or blacklisted through command line options.

Each check is independent and order is not and will not be guaranteed, and its input will be informed through options in
the command line interface; currently the only input is the required `uri` option.

The following checks have been implemented:

| Name | Description
|---|---
| `is-helm-v3` | Checks whether the given `uri` is a Helm v3 chart.
| `has-readme` | Checks whether the Helm chart contains a `README.md` file.
| `contains-test` | Checks whether the Helm chart contains at least one test file.
| `has-minkubeversion` | Checks whether the Helm chart's `Chart.yaml` includes the `minKubeVersion` field.
| `readme-contains-values-schema` | Checks whether the Helm chart `README.md` file contains a `values` schema section.
| `not-contains-crds` | Check whether the Helm chart does not include CRDs.

The following checks are being implemented and/or considered:

| Name | Description
|---|---
| `keywords-are-openshift-categories` | Checks whether the Helm chart's `Chart.yaml` file includes keywords mapped to OpenShift categories.
| `is-commercial-chart` | Checks whether the Helm chart is a Commercial chart.
| `is-community-chart` | Checks whether the Helm chart is a Community chart.
| `not-contains-infra-plugins-and-drivers` | Check whether the Helm chart does not include infra plugins and drivers (network, storage, hardware, etc)
| `can-be-installed-without-manual-prerequisites` |
| `can-be-installed-without-cluster-admin-privileges` |

## Architecture

This tool is part of a larger process that aims to certify Helm charts, and its sole responsibility is to ingest a Helm
chart URI (`file://`, `https?://`, etc)
and return either a *positive* result indicating the Helm chart has passed all checks, or a *negative* result indicating
which checks have failed and possibly propose solutions.

The application is separated in two pieces: a command line interface and a library. This is handy because the command
line interface is specific to the user interface, and the library can be generic enough to be used to, for example,
inspect Helm chart bytes in flight.

One positive aspect of the command line interface specificity is that its output can be tailored to the methods of
consumption the user expects; in other words, the command line interface can be programmed in such way it can be
represented as either *YAML* or *JSON* formats, in addition to a descriptive representation tailored to human actors.

The interpretation of what is considered a certified Helm chart depends on which checks the chart has been submitted to,
so this information should be present in the certificate as well.

Primitive functions to manipulate the Helm chart should be provided, since most checks involve inspecting the contents
of the chart itself; for example, whether a `README.md` file exists, or whether `README.md` contains the `values`'
specification, implicating in offering a cache API layer is required to avoid downloading and unpacking the charts for
each test.

## Building chart-verifier

To build `chart-verifier` locally, please execute `hack/build.sh` or its PowerShell alternative.

To build `chart-verifier` container image, please execute `hack/build-image.sh` or its PowerShell alternative:

```text
PS C:\Users\igors\GolandProjects\helmcertifier> .\hack\build-image.ps1
[+] Building 42.2s (15/15) FINISHED
 => [internal] load build definition from Dockerfile                                                                                               0.0s
 => => transferring dockerfile: 283B                                                                                                               0.0s
 => [internal] load .dockerignore                                                                                                                  0.0s
 => => transferring context: 2B                                                                                                                    0.0s
 => [internal] load metadata for docker.io/library/fedora:31                                                                                       1.4s
 => [internal] load metadata for docker.io/library/golang:1.15                                                                                     1.5s
 => [internal] load build context                                                                                                                  0.6s
 => => transferring context: 43.32MB                                                                                                               0.6s
 => CACHED [stage-1 1/2] FROM docker.io/library/fedora:31@sha256:ba4fe6a3da48addb248a16e8a63599cc5ff5250827e7232d2e3038279a0e467e                  0.0s
 => [build 1/7] FROM docker.io/library/golang:1.15@sha256:2d144ad89c91d4eb516eaf8361c1f49115c5d06042683e27d3439dc6c2535cc7                         0.0s
 => CACHED [build 2/7] WORKDIR /tmp/src                                                                                                            0.0s
 => [build 3/7] COPY go.mod .                                                                                                                      0.1s
 => [build 4/7] COPY go.sum .                                                                                                                      0.0s
 => [build 5/7] RUN go mod download                                                                                                               27.8s
 => [build 6/7] COPY . .                                                                                                                           0.1s
 => [build 7/7] RUN ./hack/build.sh                                                                                                               11.6s
 => [stage-1 2/2] COPY --from=build /tmp/src/out/chart-verifier /app/chart-verifier                                                                0.2s
 => exporting to image                                                                                                                             0.2s
 => => exporting layers                                                                                                                            0.2s
 => => writing image sha256:9a57eb6b573f3878559b44ab0a8d7be350f3ccc634b3e9b8085cbc279f3a229c                                                       0.0s
 => => naming to docker.io/library/chart-verifier:9ec6e7e                                                                                          0.0s
```

The container image created by the build program is tagged with the commit ID of the working directory at the time of
the build: `chart-verifier:9ec6e7e`.

This container image can then be executed with the Docker client as `docker run -it chart-verifier:9ec6e7e certify`,
like in the example below:

```text
PS C:\Users\igors\GolandProjects\helmcertifier> docker run -it chart-verifier:9ec6e7e certify --help
Certifies a Helm chart by checking some of its characteristics

Usage:
  chart-verifier certify [flags]

Flags:
  -e, --except strings   all available checks except those informed will be performed
  -h, --help             help for certify
  -o, --only strings     only the informed checks will be performed
  -f, --output string    the output format: default, json or yaml
  -u, --uri string       uri of the Chart being certified

Global Flags:
      --config string   config file (default is $HOME/.chart-verifier.yaml)
```

To verify a chart on the host system, the directory containing the chart should be mounted in the container; for http or
https verifications, no mounting is required:

```text
> docker run -it chart-verifier:9ec6e7e certify -u https://github.com/isutton/helmcertifier/blob/master/pk
g/chartverifier/checks/chart-0.1.0-v3.valid.tgz?raw=true
chart: chart
version: 1.16.0
ok: true

not-contains-crds:
        ok: true
        reason: Chart does not contain CRDs
helm-lint:
        ok: true
        reason: Helm lint successful
has-readme:
        ok: true
        reason: Chart has README
is-helm-v3:
        ok: true
        reason: API version is V2 used in Helm 3
contains-test:
        ok: true
        reason: Chart test files exist
contains-values:
        ok: true
        reason: Values file exist
contains-values-schema:
        ok: true
        reason: Values schema file exist
has-minkubeversion:
        ok: true
        reason: Minimum Kubernetes version specified
```

## Usage

To certify a chart against all available checks:

```text
> chart-verifier --uri ./chart.tgz
> chart-verifier --uri ~/src/chart
> chart-verifier --uri https://www.example.com/chart.tgz
```

To apply only the `is-helm-v3` check:

```text
> chart-verifier --only is-helm-v3 --uri https://www.example.com/chart.tgz
```

To apply all checks except `is-helm-v3`:

```text
> chart-verifier --except is-helm-v3 --uri https://www.example.com/chart.tgz
```
