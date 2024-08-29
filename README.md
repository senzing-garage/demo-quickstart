# demo-quickstart

If you are beginning your journey with [Senzing],
please start with [Senzing Quick Start guides].

You are in the [Senzing Garage] where projects are "tinkered" on.
Although this GitHub repository may help you understand an approach to using Senzing,
it's not considered to be "production ready" and is not considered to be part of the Senzing product.
Heck, it may not even be appropriate for your application of Senzing!

## :warning: WARNING: demo-quickstart is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

`demo-quickstart` is a command in the [senzing-tools] suite of tools.
This command creates an environment for exploring Senzing.

[![Go Reference Badge]][Package reference]
[![Go Report Card Badge]][Go Report Card]
[![License Badge]][License]
[![go-test-linux.yaml Badge]][go-test-linux.yaml]
[![go-test-darwin.yaml Badge]][go-test-darwin.yaml]
[![go-test-windows.yaml Badge]][go-test-windows.yaml]

[![golangci-lint.yaml Badge]][golangci-lint.yaml]

## Overview

`demo-quickstart` starts the Senzing gRPC server and HTTP server for use in Senzing exploration.

Senzing SDKs for accessing the gRPC server:

1. Go: [sz-sdk-go-grpc]
1. Python: [sz-sdk-python-grpc]
A simple demonstration using `senzing-tools` and a SQLite database.

```console
export LD_LIBRARY_PATH=/opt/senzing/er/lib/
export SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@/tmp/sqlite/G2C.db
senzing-tools init-database
senzing-tools demo-quickstart

```

Then visit [localhost:8261]

## Install

1. The `demo-quickstart` command is installed with the [senzing-tools] suite of tools.
   See [senzing-tools install].

## Use

```console
export LD_LIBRARY_PATH=/opt/senzing/er/lib/
senzing-tools demo-quickstart [flags]
```

1. For options and flags:
    1. [Online documentation]
    1. Runtime documentation:

        ```console
        export LD_LIBRARY_PATH=/opt/senzing/er/lib/
        senzing-tools demo-quickstart --help
        ```

1. In addition to the following simple usage examples, there are additional [Examples].

### Using command line options

1. :pencil2: Specify database using command line option.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/er/lib/
    senzing-tools demo-quickstart \
        --database-url postgresql://username:password@postgres.example.com:5432/G2 \

    ```

1. Visit [localhost:8261]
1. Run `senzing-tools demo-quickstart --help` or see [Parameters] for additional parameters.

### Using environment variables

1. :pencil2: Specify database using environment variable.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/er/lib/
    export SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2
    senzing-tools demo-quickstart
    ```

1. Visit [localhost:8261]
1. Run `senzing-tools demo-quickstart --help` or see [Parameters] for additional parameters.

### Using Docker

This usage shows how to initialize a database with a Docker container.

1. This usage specifies a URL of an external database.
   Example:

    ```console
    docker run \
        --publish 8260:8260 \
        --publish 8261:8261 \
        --rm \
        senzing/senzing-tools demo-quickstart

    ```

1. Visit [localhost:8261]
1. See [Parameters] for additional parameters.

### Parameters

- **[SENZING_TOOLS_DATABASE_URL]**
- **[SENZING_TOOLS_ENGINE_CONFIGURATION_JSON]**
- **[SENZING_TOOLS_ENGINE_LOG_LEVEL]**
- **[SENZING_TOOLS_ENGINE_MODULE_NAME]**
- **[SENZING_TOOLS_GRPC_PORT]**
- **[SENZING_TOOLS_LOG_LEVEL]**

## References

1. [Command reference]
1. [Development]
1. [Errors]
1. [Examples]
1. [DockerHub]

[Command reference]: https://garage.senzing.com/senzing-tools/senzing-tools_demo-quickstart.html
[Development]: docs/development.md
[DockerHub]: https://hub.docker.com/r/senzing/demo-quickstart
[Errors]: docs/errors.md
[Examples]: docs/examples.md
[localhost:8261]: http://localhost:8261
[Online documentation]: https://hub.senzing.com/senzing-tools/senzing-tools_demo-quickstart.html
[Parameters]: #parameters
[Senzing Garage]: https://github.com/senzing-garage-garage
[Senzing Quick Start guides]: https://docs.senzing.com/quickstart/
[SENZING_TOOLS_DATABASE_URL]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_database_url
[SENZING_TOOLS_ENGINE_CONFIGURATION_JSON]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_configuration_json
[SENZING_TOOLS_ENGINE_LOG_LEVEL]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_log_level
[SENZING_TOOLS_ENGINE_MODULE_NAME]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_module_name
[SENZING_TOOLS_GRPC_PORT]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_grpc_port
[SENZING_TOOLS_LOG_LEVEL]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_log_level
[senzing-tools install]: https://github.com/senzing-garage/senzing-tools#install
[senzing-tools]: https://github.com/senzing-garage/senzing-tools
[Senzing]: https://senzing.com/
[API documentation]: https://pkg.go.dev/github.com/senzing-garage/template-go
[Development]: docs/development.md
[DockerHub]: https://hub.docker.com/r/senzing/template-go
[Errors]: docs/errors.md
[Examples]: docs/examples.md
[Go Package library]: https://pkg.go.dev
[Go Reference Badge]: https://pkg.go.dev/badge/github.com/senzing-garage/template-go.svg
[Go Report Card Badge]: https://goreportcard.com/badge/github.com/senzing-garage/template-go
[Go Report Card]: https://goreportcard.com/report/github.com/senzing-garage/template-go
[go-test-darwin.yaml Badge]: https://github.com/senzing-garage/template-go/actions/workflows/go-test-darwin.yaml/badge.svg
[go-test-darwin.yaml]: https://github.com/senzing-garage/template-go/actions/workflows/go-test-darwin.yaml
[go-test-linux.yaml Badge]: https://github.com/senzing-garage/template-go/actions/workflows/go-test-linux.yaml/badge.svg
[go-test-linux.yaml]: https://github.com/senzing-garage/template-go/actions/workflows/go-test-linux.yaml
[go-test-windows.yaml Badge]: https://github.com/senzing-garage/template-go/actions/workflows/go-test-windows.yaml/badge.svg
[go-test-windows.yaml]: https://github.com/senzing-garage/template-go/actions/workflows/go-test-windows.yaml
[golangci-lint.yaml Badge]: https://github.com/senzing-garage/template-go/actions/workflows/golangci-lint.yaml/badge.svg
[golangci-lint.yaml]: https://github.com/senzing-garage/template-go/actions/workflows/golangci-lint.yaml
[License Badge]: https://img.shields.io/badge/License-Apache2-brightgreen.svg
[License]: https://github.com/senzing-garage/template-go/blob/main/LICENSE
[main.go]: main.go
[Package reference]: https://pkg.go.dev/github.com/senzing-garage/template-go
[Senzing Garage]: https://github.com/senzing-garage
[Senzing Quick Start guides]: https://docs.senzing.com/quickstart/
[Senzing]: https://senzing.com/