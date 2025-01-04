[![Go Reference](https://pkg.go.dev/badge/github.com/leg100/surl/v2.svg)](https://pkg.go.dev/github.com/leg100/surl/v2)
# surl

Create signed URLs in Go.

## Installation

`go get github.com/leg100/surl/v2@latest`

## Usage

```golang
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/leg100/surl/v2"
)

func main() {
	signer := surl.New([]byte("secret_key"))

	// Create a signed URL that expires in one hour.
	signed, _ := signer.Sign("https://example.com/a/b/c?foo=bar", time.Now().Add(time.Hour))
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
Benchmark/sign/path/decimal/no_opts-16            790069              1366 ns/op
Benchmark/verify/path/decimal/no_opts-16          865004              1267 ns/op
Benchmark/sign/path/decimal/prefix-16             671443              1659 ns/op
Benchmark/verify/path/decimal/prefix-16           836241              1303 ns/op
Benchmark/sign/path/decimal/skip_query-16         813769              1354 ns/op
Benchmark/verify/path/decimal/skip_query-16       869764              1262 ns/op
Benchmark/sign/path/decimal/skip_scheme-16        799594              1351 ns/op
Benchmark/verify/path/decimal/skip_scheme-16      865023              1262 ns/op
Benchmark/sign/path/decimal/prefix_and_skip_query-16              681868              1641 ns/op
Benchmark/verify/path/decimal/prefix_and_skip_query-16            835401              1288 ns/op
Benchmark/sign/path/decimal/prefix_and_skip_scheme-16             671688              1654 ns/op
Benchmark/verify/path/decimal/prefix_and_skip_scheme-16           837308              1285 ns/op
Benchmark/sign/path/decimal/prefix_and_skip_query_and_skip_scheme-16              677064              1645 ns/op
Benchmark/verify/path/decimal/prefix_and_skip_query_and_skip_scheme-16            859400              1283 ns/op
Benchmark/sign/path/base58/no_opts-16                                             832389              1332 ns/op
Benchmark/verify/path/base58/no_opts-16                                           894231              1238 ns/op
Benchmark/sign/path/base58/prefix-16                                              678120              1627 ns/op
Benchmark/verify/path/base58/prefix-16                                            832461              1280 ns/op
Benchmark/sign/path/base58/skip_query-16                                          829882              1307 ns/op
Benchmark/verify/path/base58/skip_query-16                                        903090              1218 ns/op
Benchmark/sign/path/base58/skip_scheme-16                                         785352              1322 ns/op
Benchmark/verify/path/base58/skip_scheme-16                                       886584              1233 ns/op
Benchmark/sign/path/base58/prefix_and_skip_query-16                               691098              1586 ns/op
Benchmark/verify/path/base58/prefix_and_skip_query-16                             866223              1256 ns/op
Benchmark/sign/path/base58/prefix_and_skip_scheme-16                              679152              1609 ns/op
Benchmark/verify/path/base58/prefix_and_skip_scheme-16                            872644              1274 ns/op
Benchmark/sign/path/base58/prefix_and_skip_query_and_skip_scheme-16               675469              1591 ns/op
Benchmark/verify/path/base58/prefix_and_skip_query_and_skip_scheme-16             871072              1237 ns/op
Benchmark/sign/query/decimal/no_opts-16                                           347485              3382 ns/op
Benchmark/verify/query/decimal/no_opts-16                                         367434              3140 ns/op
Benchmark/sign/query/decimal/prefix-16                                            321260              3646 ns/op
Benchmark/verify/query/decimal/prefix-16                                          361286              3153 ns/op
Benchmark/sign/query/decimal/skip_query-16                                        276272              4211 ns/op
Benchmark/verify/query/decimal/skip_query-16                                      293073              3852 ns/op
Benchmark/sign/query/decimal/skip_scheme-16                                       352858              3281 ns/op
Benchmark/verify/query/decimal/skip_scheme-16                                     373916              3035 ns/op
Benchmark/sign/query/decimal/prefix_and_skip_query-16                             267811              4353 ns/op
Benchmark/verify/query/decimal/prefix_and_skip_query-16                           301072              3875 ns/op
Benchmark/sign/query/decimal/prefix_and_skip_scheme-16                            326252              3505 ns/op
Benchmark/verify/query/decimal/prefix_and_skip_scheme-16                          370263              3059 ns/op
Benchmark/sign/query/decimal/prefix_and_skip_query_and_skip_scheme-16             267507              4341 ns/op
Benchmark/verify/query/decimal/prefix_and_skip_query_and_skip_scheme-16           298611              3857 ns/op
Benchmark/sign/query/base58/no_opts-16                                            347785              3244 ns/op
Benchmark/verify/query/base58/no_opts-16                                          376110              3032 ns/op
Benchmark/sign/query/base58/prefix-16                                             324220              3589 ns/op
Benchmark/verify/query/base58/prefix-16                                           331550              3107 ns/op
Benchmark/sign/query/base58/skip_query-16                                         291204              4002 ns/op
Benchmark/verify/query/base58/skip_query-16                                       302330              3834 ns/op
Benchmark/sign/query/base58/skip_scheme-16                                        328934              3349 ns/op
Benchmark/verify/query/base58/skip_scheme-16                                      353527              3185 ns/op
Benchmark/sign/query/base58/prefix_and_skip_query-16                              266244              4366 ns/op
Benchmark/verify/query/base58/prefix_and_skip_query-16                            300997              3846 ns/op
Benchmark/sign/query/base58/prefix_and_skip_scheme-16                             328764              3506 ns/op
Benchmark/verify/query/base58/prefix_and_skip_scheme-16                           378475              3174 ns/op
Benchmark/sign/query/base58/prefix_and_skip_query_and_skip_scheme-16              268564              4276 ns/op
Benchmark/verify/query/base58/prefix_and_skip_query_and_skip_scheme-16            295378              3939 ns/op
PASS
ok      github.com/leg100/surl  64.305s
```

