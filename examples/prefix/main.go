package main

import (
	"fmt"
	"log"
	"time"

	"github.com/leg100/surl"
)

func main() {
	signer := surl.New([]byte("secret_key"), surl.PrefixPath("/signed"))

	// Create a signed URL that expires in one hour.
	signed, _ := signer.Sign("https://example.com/a/b/c?foo=bar", time.Now().Add(time.Hour))
	fmt.Println(signed)
	// Outputs something like:
	// https://example.com/signed/a/b/c?expiry=1669574398&foo=bar&signature=NvIrIFcc1OaKgeVSN685tSD26PTdjlUxxSZRE18Wk_8

	err := signer.Verify(signed)
	if err != nil {
		log.Fatal("verification failed: ", err.Error())
	}
	fmt.Println("verification succeeded")
}
