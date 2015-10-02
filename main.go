package main

import (
	"fmt"
	"os"

	"github.com/go-microservices/presigner/option"
	"github.com/go-microservices/presigner/server"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() (err error) {
	var (
		o option.Options
	)

	o, err = option.New(os.Args[1:])
	if err != nil {
		return
	}

	err = o.ReadPrivateKey()
	if err != nil {
		return
	}

	err = server.Serve(o)
	if err != nil {
		return
	}

	return
}
