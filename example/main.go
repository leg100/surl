package main

import (
	"fmt"
	"time"

	"github.com/leg100/signer"
)

func main() {
	sign := signer.New([]byte("secret_sesame"))

	// Create a signed URL that expires in one hour.
	signed, _ := sign.Sign("https://example.com/a/b/c?foo=bar", time.Hour)

	fmt.Println("signed url:", signed)

	// Outputs something like:
	//
	// https://example.com/signed/pTn2am3eh8Ndz7ZTb6ya2gOMA5XtnFRd-1M__TNQr9o.1664441797/a/b/c?foo

	err := sign.Verify(signed)
	if err != nil {
		fmt.Println("verification failed:", err.Error())
	}
	fmt.Println("verification succeeded")
}
