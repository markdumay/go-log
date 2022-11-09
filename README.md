# go-log

<!-- Tagline -->
<p align="center">
    <b>Simplified application logger powered by zerolog</b>
    <br />
</p>


<!-- Badges -->
<p align="center">
    <a href="https://pkg.go.dev/go.markdumay.org/log" alt="Go Package">
        <img src="https://pkg.go.dev/badge/go.markdumay.org/log.svg" alt="Go Reference" />
    </a>
    <a href="https://github.com/markdumay/go-log/releases/latest" alt="Go Module">
        <img src="https://img.shields.io/github/v/tag/markdumay/go-log?label=module" alt="Go Module" />
    </a>
    <a href="https://www.codefactor.io/repository/github/markdumay/go-log" alt="CodeFactor">
        <img src="https://img.shields.io/codefactor/grade/github/markdumay/go-log" />
    </a>
    <a href="https://github.com/markdumay/go-log/commits/main" alt="Last commit">
        <img src="https://img.shields.io/github/last-commit/markdumay/go-log.svg" />
    </a>
    <a href="https://github.com/markdumay/go-log/issues" alt="Issues">
        <img src="https://img.shields.io/github/issues/markdumay/go-log.svg" />
    </a>
    <a href="https://github.com/markdumay/go-log/pulls" alt="Pulls">
        <img src="https://img.shields.io/github/issues-pr-raw/markdumay/go-log.svg" />
    </a>
    <a href="https://github.com/markdumay/go-log/blob/main/LICENSE" alt="License">
        <img src="https://img.shields.io/github/license/markdumay/go-log" />
    </a>
</p>

<!-- Table of Contents -->
<p align="center">
  <a href="#about">About</a> •
  <a href="#built-with">Built With</a> •
  <a href="#prerequisites">Prerequisites</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#contributing">Contributing</a> •
  <a href="#donate">Donate</a> •
  <a href="#license">License</a>
</p>


## About
go-log is a simplified logger package for Go applications. Using the Zero Allocation JSON Logger (zerolog) under the hood, it simplifies the logging of application-wide messages. It supports three logging modes: Default, Pretty, and JSON. Logs are directed to the console by default, but can be buffered or redirected to a log file instead.

## Built With
The project uses the following core software components:
* [Zero Allocation JSON Logger][zerolog_url] - Go package providing a fast and simple logger dedicated to JSON output.
* [Testify][testify_url] - Go unit-testing toolkit with common assertions and mocks.

## Prerequisites
go-log requires Go version 1.16 or later to be installed on your system.

## Installation
```console
go get -u go.markdumay.org/log
```

## Usage
Import go-log into your application to start using the logger. By default, go-log writes the log messages to the console. Please refer to the [package documentation][package] for more details. The following code snippet illustrates the basic usage of go-log.

```go
package main

import (
    "go.markdumay.org/log"
)

func main() {
    // show an info message using default formatting, expected output:
    // This is an info log
    log.Info("This is an info log")

    // show an error message using default formatting, expected output:
    // ERROR  Error message
    log.Info("Error message")

    // switch to pretty formatting
    log.InitLogger(log.Pretty)

    // show a warning using pretty formatting, expected output:
    // 2006-01-02T15:04:05Z07:00 | WARN   | Warning
    log.Warn("Warning")

    // switch to JSON formatting
    log.InitLogger(log.JSON)

    // switch to debug level as minimum level
    log.SetGlobalLevel(log.DebugLevel)

    // show a debug message using JSON formatting, expected output:
    // {"level":"debug","time":"2006-01-02T15:04:05Z07:00","message":"Testing level debug"}
    log.Debugf("Testing level %s", "debug")
}
```

## Contributing
go-log welcomes contributions of any kind. It is recommended to create an issue to discuss your intended contribution before submitting a larger pull request though. Please consider the following guidelines when contributing:
- Address all linting recommendations from `golangci-lint run` (using `.golangci.yml` from the repository).
- Ensure the code is covered by one or more unit tests (using [Testify][testify_url] when applicable).
- Follow the recommendations from [Effective Go][effective_go] and the [Uber Go Style Guide][uber_go_guide].

The following steps decribe how to submit a Pull Request:
1. Clone the repository and create a new branch 
    ```console
    $ git checkout https://github.com/markdumay/go-log.git -b name_for_new_branch
    ```
2. Make and test the changes
3. Submit a Pull Request with a comprehensive description of the changes

## Donate
<a href="https://www.buymeacoffee.com/markdumay" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/lato-orange.png" alt="Buy Me A Coffee" style="height: 51px !important;width: 217px !important;"></a>

## License
The go-log codebase is released under the [MIT license][license]. The documentation (including the "README") is licensed under the Creative Commons ([CC BY-NC 4.0)][cc-by-nc-4.0] license.

<!-- MARKDOWN PUBLIC LINKS -->
[cc-by-nc-4.0]: https://creativecommons.org/licenses/by-nc/4.0/
[effective_go]: https://golang.org/doc/effective_go
[testify_url]: https://github.com/stretchr/testify
[uber_go_guide]: https://github.com/uber-go/guide/
[zerolog_url]: https://github.com/rs/zerolog

<!-- MARKDOWN MAINTAINED LINKS -->
<!-- TODO: add blog link
[blog]: https://markdumay.com
-->
[blog]: https://github.com/markdumay
[license]: https://github.com/markdumay/go-log/blob/main/LICENSE
[package]: https://pkg.go.dev/go.markdumay.org/log
[repository]: https://github.com/markdumay/go-log.git
