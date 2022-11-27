[![Go Report Card](https://goreportcard.com/badge/github.com/leg100/surl)](https://goreportcard.com/report/github.com/leg100/surl)
[![Version](https://img.shields.io/badge/goversion-1.19.x-blue.svg)](https://golang.org)
[![Go Reference](https://pkg.go.dev/badge/github.com/leg100/surl.svg)](https://pkg.go.dev/github.com/leg100/surl)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/leg100/goblender/master/LICENSE)
![Tests](https://github.com/leg100/signer/actions/workflows/tests.yml/badge.svg)
# surl

Create signed URLs using go.

## Installation

`go get github.com/leg100/surl@latest`

## Usage

```golang
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/leg100/surl"
)

func main() {
	signer := surl.New([]byte("secret_sesame"))

	// Create a signed URL that expires in one hour.
	signed, _ := signer.Sign("https://example.com/a/b/c?foo=bar", time.Hour)
	fmt.Println("signed url:", signed)
	// Outputs something like:
	// https://example.com/a/b/c?expiry=1667331055&foo=bar&signature=TGvxmRwpoAUt9YEIbeJ164lMYrzA2DBnYB9Lcy9m1T

	err := signer.Verify(signed)
	if err != nil {
		log.Fatal("verification failed: ", err.Error())
	}
	fmt.Println("verification succeeded")
}
```

## Options

The format and behaviour of signed URLs can be configured by passing options to the constructor.

### Query Formatter

```go
surl.New(secret, WithQueryFormatter())
```
The query formatter is the default format. It stores the signature and expiry in query parameters.

### Path Formatter

```go
surl.New(secret, WithPathFormatter())
```

The path formatter stores the signature and expiry in the path itself:

```bash
https://example.com/PaMIbZQ6wxPdHXVLfIGwZBULo-FSTdt7-bCLZjBPPUE.1669574162/a/b/c?foo=bar
```

## Notes

* Any change in the order of the query parameters in a signed URL renders it invalid, unless `SkipQuery` is specified.
