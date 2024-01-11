# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
[markdownlint](https://dlaa.me/markdownlint/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2024-01-04

### Changed in 0.2.0

- Renamed module to `github.com/senzing-garage/demo-quickstart`
- Refactor to [template-go](https://github.com/senzing-garage/template-go)
- Update dependencies
  - github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c
  - github.com/senzing-garage/go-cmdhelping v0.2.0
  - github.com/senzing-garage/go-grpcing v0.2.0
  - github.com/senzing-garage/go-observing v0.3.0
  - github.com/senzing-garage/go-rest-api-service v0.2.0
  - github.com/spf13/cobra v1.8.0
  - github.com/spf13/viper v1.18.2
  - google.golang.org/grpc v1.60.1

## [0.1.2] - 2023-11-02

### Changed in 0.1.2

- Update dependencies
  - github.com/senzing-garage/go-rest-api-service v0.1.1

## [0.1.1] - 2023-10-25

### Changed in 0.1.1

- Refactor to [template-go](https://github.com/senzing-garage/template-go)
- Update dependencies
  - github.com/senzing-garage/go-cmdhelping v0.1.9
  - github.com/senzing-garage/go-grpcing v0.1.3
  - github.com/senzing-garage/go-observing v0.2.8
  - github.com/senzing-garage/go-rest-api-service v0.1.0
  - github.com/spf13/viper v1.17.0
  - google.golang.org/grpc v1.59.0

## [0.1.0] - 2023-10-03

### Changed in 0.1.0

- Supports SenzingAPI 3.8.0
- Deprecated functions have been removed
- Update dependencies
  - google.golang.org/grpc v1.58.2

## [0.0.5] - 2023-09-01

### Changed in 0.0.5

- Last version before SenzingAPI 3.8.0

## [0.0.4] - 2023-08-09

### Changed in 0.0.4

- Refactor to `template-go`
- Update dependencies
  - github.com/senzing-garage/go-cmdhelping v0.1.5
  - github.com/senzing-garage/go-grpcing v0.1.2
  - github.com/senzing-garage/go-observing v0.2.7
  - github.com/senzing-garage/go-rest-api-service v0.0.5

## [0.0.3] - 2023-07-26

### Changed in 0.0.3

- Update dependencies
  - github.com/senzing-garage/go-cmdhelping v0.1.2
  - github.com/senzing-garage/go-rest-api-service v0.0.4
  - google.golang.org/grpc v1.57.0

## [0.0.2] - 2023-07-21

### Added in 0.0.2

- Opens web browser, unless disabled by `--tty-only`

### Changed in 0.0.2

- Update dependencies
  - github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8
  - github.com/senzing-garage/go-cmdhelping v0.1.1
  - github.com/senzing-garage/go-common v0.2.4
  - github.com/senzing-garage/go-rest-api-service v0.0.3
  - google.golang.org/grpc v1.56.2

## [0.0.1] - 2023-06-19

### Added to 0.0.1

- Swagger UI
- Support for following Senzing HTTP REST URIs:
  - GET /heartbeat
  - GET /specifications/open-api
  - GET /license
  - GET /version
- XTerm
