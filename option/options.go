package option

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/alecthomas/kingpin"
)

type Options struct {
	GoogleAccessID string
	PrivateKeyPath string
	PrivateKey     []byte
	Buckets        Buckets
	Port           int
	Duration       time.Duration
}

func New(args []string) (o Options, err error) {
	app := kingpin.New("signing", "Publisher of signed URLs to upload files directly to Google Cloud Storage")
	i := app.Flag("id", "Google Access ID").Short('i').OverrideDefaultFromEnvar("GOOGLE_ACCESS_ID").Required().String()
	k := app.Flag("key", "Path to private key").Short('k').OverrideDefaultFromEnvar("PRIVATE_KEY_PATH").Required().String()
	b := app.Flag("buckets", "Allowed buckets").Short('b').Default("*").OverrideDefaultFromEnvar("BUCKETS").String()
	p := app.Flag("port", "Listening port").Short('p').Default("80").OverrideDefaultFromEnvar("PORT").Int()
	d := app.Flag("duration", "Available duration of published signature").Short('d').Default("1m").Duration()
	_, err = app.Parse(args)
	if err != nil {
		return
	}

	o.GoogleAccessID = *i
	o.PrivateKeyPath = *k
	o.Buckets = NewBuckets(*b)
	o.Port = *p
	o.Duration = *d
	return
}

func (o *Options) ReadPrivateKey() (err error) {
	privateKey, err := ioutil.ReadFile(o.PrivateKeyPath)
	if err != nil {
		return
	}
	o.PrivateKey = privateKey
	return
}

type Buckets struct {
	wildcard  bool
	whitelist []string
}

func NewBuckets(buckets string) (b Buckets) {
	b.whitelist = strings.Split(buckets, ",")
	for _, white := range b.whitelist {
		if white == "*" {
			b.wildcard = true
			break
		}
	}
	return
}

func (b Buckets) Contains(bucket string) bool {
	if b.wildcard {
		return true
	}
	for _, white := range b.whitelist {
		if white == bucket {
			return true
		}
	}
	return false
}
