[![codecov](https://codecov.io/gh/michaeljsaenz/kui/branch/main/graph/badge.svg?token=FF4ZXBZCBC)](https://codecov.io/gh/michaeljsaenz/kui)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmichaeljsaenz%2Fkui.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fmichaeljsaenz%2Fkui?ref=badge_shield)

# KUI
GUI for kubectl


## TODO
- [ ] add test coverage
- [ ] catch panic when cluster context not available: `panic: Get "https://1.2.3.4:443/api/v1/pods": dial tcp 1.2.3.4:443: i/o timeout`
- [ ] test for kubeConfig not accessible or set
- [ ] test for clusterContext not set or empty
- [ ] clear podTab data on refresh (similar to podLogTab data on refresh)
- [ ] optimize the log tabs (load time)

## Features
- [ ]  add copy capability to UI

## Release
1.0.0 (TBD)

- [ ] add initial build release (run `make build` to build binary and package app)

## Run tests locally
```
go test -race -coverprofile=coverage.txt -covermode=atomic
```
