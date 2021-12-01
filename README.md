# gojen

[![Release](https://github.com/Hunter-Thompson/gojen/actions/workflows/release.yml/badge.svg)](https://github.com/Hunter-Thompson/gojen/actions/workflows/release.yml) [![Upload Binary](https://github.com/Hunter-Thompson/gojen/actions/workflows/upload-binary.yml/badge.svg)](https://github.com/Hunter-Thompson/gojen/actions/workflows/upload-binary.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/Hunter-Thompson/gojen)](https://goreportcard.com/report/github.com/Hunter-Thompson/gojen) [![codecov](https://codecov.io/gh/Hunter-Thompson/gojen/branch/master/graph/badge.svg?token=LC2KR2180N)](https://codecov.io/gh/Hunter-Thompson/gojen)  

Define your go project's configuration using a json config. This config can be used to generate a new go project for you, and can also create configs for different tools.

## Features

- Release workflow

If you choose to enable the release workflow, gojen creates a Github actions workflow that simply runs `gojen`, which lints your code using `golangci-lint`, runs your tests according to the config, and then creates a release on Github using `semantic-go-release`. After that it uploads a binary of your built project as an attachment.

- Build workflow

gojen creates a Github workflow that runs the `gojen` cli, which builds lints/tests/builds your project.

- Gitignore

You can define all your gitignore entries inside `gojen.json`

- Codeowners 

You can define all your codeowners entries inside `gojen.json`

- Golangci-lint

Lint your code using golangci-lint

- Gotest 

Arguments for the go test command

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
  ]
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

## Known issues

1. The default GITHUB_TOKEN cannot be used for the release workflow since if the release is created by the github bot, the upload binary workflow will not run.

So please upload a GITHUB token inside the repo secrets, and add the key to the secret inside `gojen.json`

## Example project

https://github.com/Hunter-Thompson/test-gojen

https://github.com/Hunter-Thompson/gojen - yes gojen uses gojen : )

## Notes

This project was inspired by #projen/projen.



