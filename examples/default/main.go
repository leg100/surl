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
