[![Go Report Card](https://goreportcard.com/badge/github.com/tsawler/signer)](https://goreportcard.com/report/github.com/tsawler/signer)
[![Version](https://img.shields.io/badge/goversion-1.18.x-blue.svg)](https://golang.org)
<a href="https://golang.org"><img src="https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat-square" alt="Built with GoLang"></a>
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/tsawler/goblender/master/LICENSE)
![Tests](https://github.com/tsawler/signer/actions/workflows/tests.yml/badge.svg)
<a href="#!"><img src="https://img.shields.io/badge/coverage-100%25-brightgreen"></a>
# Signer

Signer is a simple package that makes signing URLs painless. It uses
[github.com/bwmarrin/go-alone](https://github.com/bwmarrin/go-alone) to sign URLs.

This is useful for things like sending an email with a link that can be verified, and which is
tamper-proof.

## Installation

`go get github.com/tsawler/signer@latest`

## Usage

```golang

package main

import (
	"fmt"
	"github.com/tsawler/signer"
)

const secret = "somelongsecuresecret"

func main() {
	// Create a variable of type Signature, and pass it a secret, <= 64 characters.
	sign := signer.Signature{Secret: secret}

	// Call the SignURL to get a signed version. Note that only the part after 
	// https://somesite.com or http://somesite.com is actually signed, but you 
	// must pass the full url. This way, we can use the package in development 
	// without worrying about the domain name of a particular site.
	signed, _ := sign.SignURL("https://example.com/test?id=1")
	fmt.Println("Signed url:", signed)
	
	// Output is something like:
	// https://example.com/test?id=1&hash=.3w4TgJ.pAJWBPAO5k1cimZJ-nrRKnlvosOY1Krrp3ALf1rOAds
	
	// Verify that a signed URL is valid, and was  issued by this application. Here, 
	// valid is true if the URL has a valid signature, and false if it is not.
	valid, _ := sign.VerifyURL(signed)
	fmt.Println("Valid url:", valid)

	// You can also check for expiry. Here, the signed url expires after 30 minutes.
	expired := sign.Expired(signed, 30)
	fmt.Println("Expired:", expired)
}
```