[![Go Report Card](https://goreportcard.com/badge/github.com/leg100/surl)](https://goreportcard.com/report/github.com/leg100/surl)
[![Version](https://img.shields.io/badge/goversion-1.19.x-blue.svg)](https://golang.org)
[![Go Reference](https://pkg.go.dev/badge/github.com/leg100/surl.svg)](https://pkg.go.dev/github.com/leg100/surl)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/leg100/goblender/master/LICENSE)
![Tests](https://github.com/leg100/signer/actions/workflows/tests.yml/badge.svg)
# surl

Create signed URLs in Go.

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
	signer := surl.New([]byte("secret_key"))

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

#### Query Formatter

```go
surl.New(secret, surl.WithQueryFormatter())
```
The query formatter is the default format. It stores the signature and expiry in query parameters:

```bash
https://example.com/a/b/c?expiry=1667331055&foo=bar&signature=TGvxmRwpoAUt9YEIbeJ164lMYrzA2DBnYB9Lcy9m1T
```

#### Path Formatter

```go
surl.New(secret, surl.WithPathFormatter())
```

The path formatter stores the signature and expiry in the path itself:

```bash
https://example.com/PaMIbZQ6wxPdHXVLfIGwZBULo-FSTdt7-bCLZjBPPUE.1669574162/a/b/c?foo=bar
```

#### Prefix Path

```go
surl.New(secret, surl.PrefixPath("/signed"))
```

Prefix the signed URL path:

```bash
https://example.com/signed/a/b/c?expiry=1669574398&foo=bar&signature=NvIrIFcc1OaKgeVSN685tSD26PTdjlUxxSZRE18Wk_8
```

Note: a slash is implicitly inserted between the prefix and the rest of the path.

#### Skip Query

```go
surl.New(secret, surl.SkipQuery())
```

Skip the query string when computing the signature. This is useful, say, if you have pagination query parameters but you want to use the same signed URL regardless of their value. See the [example](./examples/skip_query/main.go).

#### Skip Scheme

```go
surl.New(secret, surl.SkipScheme())
```

Skip the scheme when computing the signature. This is useful, say, if you generate signed URLs in production where you use https but you want to use these URLs in development too where you use http. See the [example](./examples/skip_scheme/main.go).

#### Decimal Encoding of Expiry

```go
surl.New(secret, surl.WithDecimalExpiry())
```

Encode expiry in decimal. This is the default.

```bash
https://example.com/a/b/c?expiry=1667331055&foo=bar&signature=TGvxmRwpoAUt9YEIbeJ164lMYrzA2DBnYB9Lcy9m1T
```

#### Base58 Encoding of Expiry

```go
surl.New(secret, surl.WithBase58Expiry())
```

Encode the expiry using Base58:

```bash
https://example.com/a/b/c?expiry=3xx1vi&foo=bar&signature=-mwCtMLTBgDkShZTbBcHjRCRXtO_ZYPE0cmrh3u6S-s
```

## Notes

* Any change in the order of the query parameters in a signed URL renders it invalid, unless `SkipQuery` is specified.

## Benchmarks

```bash
> go test -run=XXX -bench=.
goos: linux
goarch: amd64
pkg: github.com/leg100/surl
cpu: AMD Ryzen 7 3800X 8-Core Processor
Benchmark/sign/path/decimal/no_opts-16            460861              2421 ns/op
Benchmark/verify/path/decimal/no_opts-16          669487              1820 ns/op
Benchmark/sign/path/decimal/prefix-16             351980              2937 ns/op
Benchmark/verify/path/decimal/prefix-16           566366              1892 ns/op
Benchmark/sign/path/decimal/skip_query-16         572380              2411 ns/op
Benchmark/verify/path/decimal/skip_query-16       587379              1869 ns/op
Benchmark/sign/path/decimal/prefix_and_skip_query-16              353696              3019 ns/op
Benchmark/verify/path/decimal/prefix_and_skip_query-16            601311              1891 ns/op
Benchmark/sign/path/base58/no_opts-16                             555777              2418 ns/op
Benchmark/verify/path/base58/no_opts-16                           629716              1792 ns/op
Benchmark/sign/path/base58/prefix-16                              431458              2776 ns/op
Benchmark/verify/path/base58/prefix-16                            714262              1848 ns/op
Benchmark/sign/path/base58/skip_query-16                          448022              2380 ns/op
Benchmark/verify/path/base58/skip_query-16                        618013              1837 ns/op
Benchmark/sign/path/base58/prefix_and_skip_query-16               396823              2759 ns/op
Benchmark/verify/path/base58/prefix_and_skip_query-16             547171              1840 ns/op
Benchmark/sign/query/decimal/no_opts-16                           195300              6263 ns/op
Benchmark/verify/query/decimal/no_opts-16                         229317              5367 ns/op
Benchmark/sign/query/decimal/prefix-16                            184351              6387 ns/op
Benchmark/verify/query/decimal/prefix-16                          231315              5657 ns/op
Benchmark/sign/query/decimal/skip_query-16                        144122              7567 ns/op
Benchmark/verify/query/decimal/skip_query-16                      166753              6878 ns/op
Benchmark/sign/query/decimal/prefix_and_skip_query-16             157466              7757 ns/op
Benchmark/verify/query/decimal/prefix_and_skip_query-16           157782              6751 ns/op
Benchmark/sign/query/base58/no_opts-16                            181545              6106 ns/op
Benchmark/verify/query/base58/no_opts-16                          249518              5178 ns/op
Benchmark/sign/query/base58/prefix-16                             192388              6307 ns/op
Benchmark/verify/query/base58/prefix-16                           222874              5344 ns/op
Benchmark/sign/query/base58/skip_query-16                         177720              7256 ns/op
Benchmark/verify/query/base58/skip_query-16                       190534              6436 ns/op
Benchmark/sign/query/base58/prefix_and_skip_query-16              147140              7465 ns/op
Benchmark/verify/query/base58/prefix_and_skip_query-16            175618              6741 ns/op
PASS
ok      github.com/leg100/surl  39.958s
```

