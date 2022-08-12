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

| Name        | Values           | Description |
|-------------|------------------|-------------|
| package     | oa3gen (default) | Override the package name for the generated files
| timetype    | time   (default) | `{type: string, format: date/datetime/time}` uses time.Time
|             | chrono           | `{type: string, format: date/datetime/time}` uses chrono.X
| decimaltype | string (default) | `{type: string, format: decimal}` changes nothing
| decimaltype | shopspring       | `{type: string, format: decimal}` uses shopspring decimal
| uuidtype    | string (default) | `{type: string, format: uuid}` uses string
| uuidtype    | google           | `{type: string, format: uuid}` uses google's uuid library

### Go Client

| Name        | Values           | Description |
|-------------|------------------|-------------|
| package     | oa3gen (default) | Override the package name for the generated files
| timetype    | time   (default) | `{type: string, format: date/datetime/time}` uses time.Time
|             | chrono           | `{type: string, format: date/datetime/time}` uses chrono.X
| decimaltype | string (default) | `{type: string, format: decimal}` changes nothing
| decimaltype | shopspring       | `{type: string, format: decimal}` uses shopspring decimal
| uuidtype    | string (default) | `{type: string, format: uuid}` uses string
| uuidtype    | google           | `{type: string, format: uuid}` uses google's uuid library

### Typescript Client

Currently no parameters are supported.
