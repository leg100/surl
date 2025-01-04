package main

import (
	"fmt"
	"log"
	"time"

	"github.com/leg100/surl/v2"
)

func main() {
	signer := surl.New([]byte("secret_key"), surl.WithBase58Expiry())

	// Create a signed URL that expires in one hour.
	signed, _ := signer.Sign("https://example.com/a/b/c?foo=bar", time.Now().Add(time.Hour))
	fmt.Println(signed)
	// Outputs something like:
	// https://example.com/a/b/c?expiry=3xx1vi&foo=bar&signature=-mwCtMLTBgDkShZTbBcHjRCRXtO_ZYPE0cmrh3u6S-s

	err := signer.Verify(signed)
	if err != nil {
		log.Fatal("verification failed: ", err.Error())
	}
	fmt.Println("verification succeeded")
}
