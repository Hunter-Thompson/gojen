# gojen

[![Release](https://github.com/Hunter-Thompson/gojen/actions/workflows/release.yml/badge.svg)](https://github.com/Hunter-Thompson/gojen/actions/workflows/release.yml) [![Upload Binary](https://github.com/Hunter-Thompson/gojen/actions/workflows/upload-binary.yml/badge.svg)](https://github.com/Hunter-Thompson/gojen/actions/workflows/upload-binary.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/Hunter-Thompson/gojen)](https://goreportcard.com/report/github.com/Hunter-Thompson/gojen) [![codecov](https://codecov.io/gh/Hunter-Thompson/gojen/branch/master/graph/badge.svg?token=LC2KR2180N)](https://codecov.io/gh/Hunter-Thompson/gojen)  

Define your go project's configuration using a json config. This config can be used to generate a new go project for you, and can also create configs for different tools.

This project also helps you maintain your code after the scaffolding process, simply run `gojen` to vendor deps, tidy them, fmt your code, `golangci-lint run`, `go test` and `go build` your code.


```
$ gojen
ℹ  | Setup | running go mod vendor
ℹ  | Setup | running go mod tidy
ℹ  | Setup | running go fmt
ℹ  | Lint | running go linter
✅  | Lint | go linter passed
ℹ  | Test | running go test
✅  | Test | go test passed
ℹ  | Build | running go build
✅  | Build | go build passed
```
---
## Features

**Release workflow**

Github actions workflow that simply runs the `golangci-lint` action, followed by `gojen` to vendor/tidy/fmt/test/build our code, a release job that creates a github release, and finally another workflow that uploads our built binary to the release we just created.

Is triggered on a push to the `defaultReleaseBranch` inside `gojen.json`.

**Build workflow**

Github workflow that runs the `golangci-lint` action, followed by `gojen`. Is triggered on pull request creation.

**Gitignore**

Define all your gitignore entries inside `gojen.json`

**Codeowners**

Define all your codeowners inside `gojen.json`

**Golangci-lint**

Lint your code using golangci-lint

**go test**

Test your code using `go test`. You can also append test arguments to `go test` by adding your arguments to the `goTestArgs` slice inside `gojen.json`

**go build**

Build your binary using the `go build` command. You can also append build arguments to `go build` by adding your arguments to the `goBuildArgs` slice inside `gojen.json`


## Getting started

Install binary

```
go install github.com/hunter-thompson/gojen
```

Create a config for your project

```
printf "{
  "name": "gojen",
  "description": "Go project generator",
  "repository": "github.com/Hunter-Thompson/gojen",
  "goVersion": "1.17",
  "authorName": "Hunter Thompson",
  "authorEmail": "hunter@example.com",
  "authorOrganization": "Hunter-Thompson",
  "licensed": true,
  "readme": true,
  "release": true,
  "buildWorkflow": true,
  "githubToken": "GIT_TOKEN",
  "defaultReleaseBranch": "master",
  "isGojen": false,
  "gitignore": [
	  ".vscode",
	  ".idea",
  ],
  "codeOwners": [
	  "* Hunter-Thompson",
  ],
  "goLinter": true,
  "goTest": true,
  "codeCov": true,
  "goTestArgs": [
	  "-v",
	  "-cover",
	  "./..."
  ],
  "goBuildArgs: ["arg1"],
}" > gojen.json
```

Generate project

```
gojen new
```

Runing `gojen` after the project has been created does the following things:

- go mod vendor
- go mod tidy
- go fmt
- golangci-lint
- go test
- go build

## Known issues

1. The default GITHUB_TOKEN cannot be used for the release workflow since if the release is created by the github bot, the upload binary workflow will not run.

So please upload a GITHUB token inside the repo secrets, and add the key to the secret inside `gojen.json`

## Example project

https://github.com/Hunter-Thompson/test-gojen 
https://github.com/Hunter-Thompson/gojen - yes gojen uses gojen : )

## Notes

This project was inspired by #projen/projen.



