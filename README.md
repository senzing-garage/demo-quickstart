# playground

If you are beginning your journey with [Senzing],
please start with [Senzing Quick Start guides].

You are in the [Senzing Garage] where projects are "tinkered" on.
Although this GitHub repository may help you understand an approach to using Senzing,
it's not considered to be "production ready" and is not considered to be part of the Senzing product.
Heck, it may not even be appropriate for your application of Senzing!

## :warning: WARNING: playground is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

`playground` creates an environment for exploring Senzing.

[![Go Reference Badge]][Package reference]
[![Go Report Card Badge]][Go Report Card]
[![License Badge]][License]
[![go-test-linux.yaml Badge]][go-test-linux.yaml]
[![go-test-darwin.yaml Badge]][go-test-darwin.yaml]
[![go-test-windows.yaml Badge]][go-test-windows.yaml]

[![golangci-lint.yaml Badge]][golangci-lint.yaml]

## Overview

`playground` starts the Senzing gRPC server and HTTP server for use in Senzing exploration.
A user visits the web interface to help with the exploration.

## Install

To start the Senzing Playground, run

```console
docker run -it --name senzing-playground -p 8260:8260 -p 8261:8261 --rm senzing/playground
```

then visit [localhost:8260].

## Use

Senzing SDKs for accessing the gRPC server:

1. Go: [sz-sdk-go-grpc]
1. Python: [sz-sdk-python-grpc]

## References

1. [Command reference]
1. [Development]
1. [Errors]
1. [Examples]
1. [DockerHub]

[Command reference]: https://garage.senzing.com/senzing-tools/senzing-tools_playground.html
[Development]: docs/development.md
[DockerHub]: https://hub.docker.com/r/senzing/playground
[Errors]: docs/errors.md
[Examples]: docs/examples.md
[Go Reference Badge]: https://pkg.go.dev/badge/github.com/senzing-garage/playground.svg
[Go Report Card Badge]: https://goreportcard.com/badge/github.com/senzing-garage/playground
[Go Report Card]: https://goreportcard.com/report/github.com/senzing-garage/playground
[go-test-darwin.yaml Badge]: https://github.com/senzing-garage/playground/actions/workflows/go-test-darwin.yaml/badge.svg
[go-test-darwin.yaml]: https://github.com/senzing-garage/playground/actions/workflows/go-test-darwin.yaml
[go-test-linux.yaml Badge]: https://github.com/senzing-garage/playground/actions/workflows/go-test-linux.yaml/badge.svg
[go-test-linux.yaml]: https://github.com/senzing-garage/playground/actions/workflows/go-test-linux.yaml
[go-test-windows.yaml Badge]: https://github.com/senzing-garage/playground/actions/workflows/go-test-windows.yaml/badge.svg
[go-test-windows.yaml]: https://github.com/senzing-garage/playground/actions/workflows/go-test-windows.yaml
[golangci-lint.yaml Badge]: https://github.com/senzing-garage/playground/actions/workflows/golangci-lint.yaml/badge.svg
[golangci-lint.yaml]: https://github.com/senzing-garage/playground/actions/workflows/golangci-lint.yaml
[License Badge]: https://img.shields.io/badge/License-Apache2-brightgreen.svg
[License]: https://github.com/senzing-garage/playground/blob/main/LICENSE
[localhost:8260]: http://localhost:8260
[Package reference]: https://pkg.go.dev/github.com/senzing-garage/playground
[Senzing Garage]: https://github.com/senzing-garage-garage
[Senzing Quick Start guides]: https://docs.senzing.com/quickstart/
[Senzing]: https://senzing.com/
[sz-sdk-go-grpc]: https://github.com/senzing-garage/sz-sdk-go-grpc
[sz-sdk-python-grpc]: https://github.com/senzing-garage/sz-sdk-python-grpc
