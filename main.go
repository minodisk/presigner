package main

import (
	"fmt"
	"os"

	"github.com/minodisk/presigner/options"
	"github.com/minodisk/presigner/server"
)

func main() {
	if err := _main(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func _main() error {
	o, err := options.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	return server.Serve(o)
}
