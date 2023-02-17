<p align="center">
  <a href="https://goreportcard.com/report/github.com/michaeljsaenz/kview"><img src="https://goreportcard.com/badge/github.com/michaeljsaenz/kview" alt="Code Status" /></a>

  <a href="https://codecov.io/gh/michaeljsaenz/kview"><img src="https://codecov.io/gh/michaeljsaenz/kview/branch/main/graph/badge.svg?token=FF4ZXBZCBC" alt='Coverage Status' /></a>

  <a href="https://app.fossa.com/projects/git%2Bgithub.com%2Fmichaeljsaenz%2Fkview?ref=badge_shield" alt="FOSSA Status"><img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmichaeljsaenz%2Fkview.svg?type=shield"/></a>
</p>

# KView
UI for kubectl


## TODO
- [ ] +test coverage
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
