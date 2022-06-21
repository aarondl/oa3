# oa3

**oa3** generates HTTP endpoint stubs and supporting models from OpenAPI 3 specs.

## Requirements

* Golang >=1.14
* `git`
* `sed` - GNU sed

> **NOTE for Mac OS nerds**: BSD sed will not work. `brew install gnu-sed` then update your path. [Click here](https://stackoverflow.com/a/34815955) for more info.

## Getting Started

1. `./build.sh`
1. `go test`
1. `./oa3 --help`

## Examples

Coming Soon™

## Parameters

This is an exhaustive list of things that are supported as a paremeter

### Go Server

| Name     | Values         | Description |
|----------|----------------|-------------|
| timetype | time (default) | Uses time.Time for date/datetime/time
|          | chrono         | Uses chrono.X for date/datetime/time

### Typescript Client

Currently no parameters are supported.
