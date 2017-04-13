package options

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type Options struct {
	GoogleAccessID string
	PrivateKeyPath string
	Buckets        Buckets
	Port           int
	Duration       time.Duration
}

func New(args []string) (Options, error) {
	var o Options
	o.Buckets = Buckets{}
	fs := flag.NewFlagSet("presigner", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  presigner [options]\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
	}
	fs.StringVar(&o.GoogleAccessID, "id", "", "Google Access ID")
	fs.StringVar(&o.PrivateKeyPath, "key", "/secret/google-auth.json", "Path to private key")
	fs.Var(&o.Buckets, "bucket", "allowed buckets")
	fs.IntVar(&o.Port, "port", 80, "listening port")
	fs.DurationVar(&o.Duration, "duration", time.Minute, "Available duration of published signature")
	return o, fs.Parse(args)
}
