[![codecov](https://codecov.io/gh/michaeljsaenz/kview/branch/main/graph/badge.svg?token=FF4ZXBZCBC)](https://codecov.io/gh/michaeljsaenz/kview)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmichaeljsaenz%2Fkview.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fmichaeljsaenz%2Fkview?ref=badge_shield)

# KView
UI for kubectl


## TODO
- [ ] add test coverage
- [ ] refactor `main.go`


## Features
- [ ]  add copy capability to UI
- [ ]  add progress bar while logs load
- [ ]  dynamic list load (pod list data) content when cluster context changes
- [ ]  update `Age` to modify for days
- [ ]  retrieve lastest logs with buffer (size limit)

## Release
0.0.1 (TBD)

- [ ] add initial build release (run `make build` to build binary and package app)

## Run tests locally
```
go test -race -coverprofile=coverage.txt -covermode=atomic
```
