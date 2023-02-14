[![codecov](https://codecov.io/gh/michaeljsaenz/kui/branch/main/graph/badge.svg?token=FF4ZXBZCBC)](https://codecov.io/gh/michaeljsaenz/kui)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmichaeljsaenz%2Fkui.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fmichaeljsaenz%2Fkui?ref=badge_shield)

# KUI
GUI for kubectl


## TODO
- [ ] add test coverage
- [ ] refactor `main.go`


## Features
- [ ]  add copy capability to UI
- [ ]  add progress bar while logs load
- [ ] dynamic list load (pod list data) content when cluster context changes

## Release
1.0.0 (TBD)

- [ ] add initial build release (run `make build` to build binary and package app)

## Run tests locally
```
go test -race -coverprofile=coverage.txt -covermode=atomic
```
