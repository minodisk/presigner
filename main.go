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
		os.Exit(2)
	}
}

func _main() error {
	o := &options.Options{}
	if err := o.Parse(os.Args[1:]); err != nil {
		return err
	}
	return server.Start(o)
}
