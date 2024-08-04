# Private GoLang Module

This private module is designed for the company's Go projects, providing a set of common methods and utilities that can be helpful and convenient in development. It simplifies the process of working on different projects and ensures standard ways of solving tasks, such as logging, metrics collection, and more.

## Key Features

- **Logging**: Tools for easy and efficient logging across different log levels (debug, info, warning, and error).
- **Metrics**: Collection and aggregation of metrics for monitoring and optimizing performance.

## Installation

```
go get github.com/wavix/go-lib
```

## Usage

The example below shows how to use the module for logging:

```go
import (
  "github.com/wavix/go-lib"
)

func main() {
 	logs := logger.New("Service name", nil)
	logs.Info(&logger.LoggerEventParams{ID: "<operation id>"}).Msg("Log example")

    err := errors.New("some error")
    logs.error(nil).Msgf("Err: %v", err)
}
```
