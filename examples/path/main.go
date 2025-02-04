package main

import (
	"fmt"
	"log"
	"time"

	"github.com/leg100/surl/v2"
)

func main() {
	signer := surl.New([]byte("secret_key"), surl.WithPathFormatter())

	// Create a signed URL that expires in one hour.
	signed, _ := signer.Sign("https://example.com/a/b/c?foo=bar", time.Now().Add(time.Hour))
	fmt.Println(signed)
	// Outputs something like:
	// https://example.com/PaMIbZQ6wxPdHXVLfIGwZBULo-FSTdt7-bCLZjBPPUE.1669574162/a/b/c?foo=bar

	err := signer.Verify(signed)
	if err != nil {
		log.Fatal("verification failed: ", err.Error())
	}
	fmt.Println("verification succeeded")
}
