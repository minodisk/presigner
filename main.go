package main

import (
	"log"
	"os"

	"github.com/go-microservices/signing/option"
	"github.com/go-microservices/signing/server"
)

func main() {
	var (
		o   option.Options
		err error
	)

	o, err = option.New(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}
	err = o.ReadPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	err = server.Serve(o)
	if err != nil {
		log.Fatal(err)
	}
}
