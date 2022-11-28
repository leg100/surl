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
	fmt.Println(signed)
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
surl.New(secret, surl.WithQueryFormatter())
```
The query formatter is the default format. It stores the signature and expiry in query parameters:

```bash
https://example.com/a/b/c?expiry=1667331055&foo=bar&signature=TGvxmRwpoAUt9YEIbeJ164lMYrzA2DBnYB9Lcy9m1T
```

### Path Formatter

```go
surl.New(secret, surl.WithPathFormatter())
```

The path formatter stores the signature and expiry in the path itself:

```bash
https://example.com/PaMIbZQ6wxPdHXVLfIGwZBULo-FSTdt7-bCLZjBPPUE.1669574162/a/b/c?foo=bar
```

### Prefix Path

```go
surl.New(secret, surl.PrefixPath("/signed"))
```

Prefix the signed URL path:

```bash
https://example.com/signed/a/b/c?expiry=1669574398&foo=bar&signature=NvIrIFcc1OaKgeVSN685tSD26PTdjlUxxSZRE18Wk_8
```

Note: a slash is implicitly inserted between the prefix and the rest of the path.

### Skip Query

```go
surl.New(secret, surl.SkipQuery())
```

Skip the query string when computing the signature. This is useful, say, if you have pagination query parameters but you want to use the same signed URL regardless of their value. See the [example](./example/skip_query/main.go).

### Decimal Encoding of Expiry

```go
surl.New(secret, surl.WithDecimalExpiry())
```

Encode expiry in decimal. This is the default.

```bash
https://example.com/a/b/c?expiry=1667331055&foo=bar&signature=TGvxmRwpoAUt9YEIbeJ164lMYrzA2DBnYB9Lcy9m1T
```

### Base58 Encoding of Expiry

```go
surl.New(secret, surl.WithBase58Expiry())
```

Encode the expiry using Base58:

```bash
https://example.com/a/b/c?expiry=3xx1vi&foo=bar&signature=-mwCtMLTBgDkShZTbBcHjRCRXtO_ZYPE0cmrh3u6S-s
```

## Notes

* Any change in the order of the query parameters in a signed URL renders it invalid, unless `SkipQuery` is specified.
