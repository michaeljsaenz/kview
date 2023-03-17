<p align="center">
  <a href="https://goreportcard.com/report/github.com/michaeljsaenz/kview"><img src="https://goreportcard.com/badge/github.com/michaeljsaenz/kview" alt="Code Status" ></a>
  <a href="https://codecov.io/gh/michaeljsaenz/kview"><img src="https://codecov.io/gh/michaeljsaenz/kview/branch/main/graph/badge.svg?token=FF4ZXBZCBC" alt='Coverage Status' /></a>
  <a href="https://app.fossa.com/projects/git%2Bgithub.com%2Fmichaeljsaenz%2Fkview?ref=badge_shield" alt="FOSSA Status"><img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmichaeljsaenz%2Fkview.svg?type=shield"/></a>
  <a href="https://img.shields.io/github/v/release/michaeljsaenz/kview?include_prereleases" title="Latest Release" rel="nofollow"><img src="https://img.shields.io/github/v/release/michaeljsaenz/kview?include_prereleases" alt="Latest Release"></a>
</p>

# KView
[KView](https://kview.app) is a standalone desktop application to interact with your Kubernetes cluster.  Get started by downloading it from [KView website](https://kview.app).
- Utilizes the [Fyne toolkit](https://fyne.io/)
- Written :100: percent in [Go](https://go.dev/)
- Built for macOS:apple:
- Contributions welcome:exclamation:
  - KView source code is available to everyone under the [MIT license](./LICENSE).

## Features
- [x]  view pod logs, annotations, labels, status, events, cluster context
- [x]  search application(pod) list
- [x]  on-demand refresh


## Screenshots


## TODO
- [ ]  add copy capability to UI
- [ ]  add `describe` button or similar to pull add'l details
- [ ]  add `yaml` button, to export via `yaml`
- [ ]  add progress bar during log loading
- [ ]  refresh/update pod list data when cluster context changes
- [ ]  add [CompletionEntry](https://github.com/fyne-io/fyne-x#completionentry) to search input
- [ ]  add `tested with` k8s versions

## Release <a href="https://img.shields.io/github/v/release/michaeljsaenz/kview?include_prereleases" title="Latest Release" rel="nofollow"><img src="https://img.shields.io/github/v/release/michaeljsaenz/kview?include_prereleases" alt="Latest Release"></a>

run `make build` to build binary and package app

## Run tests locally
```
go test -race -coverprofile=coverage.txt -covermode=atomic
```
