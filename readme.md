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
	"time"

	"github.com/leg100/surl"
)

func main() {
	signer := surl.New([]byte("secret_sesame"))

	// Create a signed URL that expires in one hour.
	signed, _ := signer.Sign("https://example.com/a/b/c?foo=bar", time.Hour)
	fmt.Println("signed url:", signed)
	// Outputs something like:
	// https://example.com/signed/pTn2am3eh8Ndz7ZTb6ya2gOMA5XtnFRd-1M__TNQr9o.1664441797/a/b/c?foo=bar

	err := signer.Verify(signed)
	if err != nil {
		fmt.Println("verification failed:", err.Error())
	}
	fmt.Println("verification succeeded")
}
```


## Notes

* Only the path and query are signed; the scheme and hostname are skipped when producing the signature. The query too can be skipped with the `SkipQuery` option.
* Any change in the order of the query parameters in a signed URL renders it invalid, unless `SkipQuery` is specified.

## TODO:

* base58 encode expiry
