# oa3

**oa3** generates HTTP endpoint stubs and supporting models from OpenAPI 3 specs.

## Requirements

* Golang >=1.18
* `git`

## Getting Started using oa3

1. `go install github.com/aarondl/oa3@latest`
2. `./oa3 <your_favorite_generator>`

## Getting Started developing on oa3

1. `go build`
1. `go test`
1. `./oa3 --help`

## Examples

Coming Soon™. Currently as there are no examples you can look at the testdata
directories of your favorite language where generally you will find a test
openapi yaml spec as well as the generated files for that yaml.

## Parameters

This is an exhaustive list of things that are supported as a paremeter by each
generator.

### Go Server

| Name     | Values           | Description |
|----------|------------------|-------------|
| package  | oa3gen (default) | Override the package name for the generated files
| timetype | time (default)   | Uses time.Time for date/datetime/time
|          | chrono           | Uses chrono.X for date/datetime/time

### Go Client

| Name     | Values         | Description |
|----------|----------------|-------------|
| package  | oa3gen (default) | Override the package name for the generated files
| timetype | time (default) | Uses time.Time for date/datetime/time
|          | chrono         | Uses chrono.X for date/datetime/time

### Typescript Client

Currently no parameters are supported.
