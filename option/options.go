package option

import "github.com/alecthomas/kingpin"

type Options struct {
	GoogleAccessID string
	PrivateKeyPath string
	Port           int
}

func New(args []string) (o Options, err error) {
	app := kingpin.New("signing-gcs", "Publisher of singed form data to upload files to Amazon S3")
	i := app.Flag("id", "Google Access ID").Short('i').OverrideDefaultFromEnvar("GOOGLE_ACCESS_ID").Required().String()
	k := app.Flag("key", "Path to private key").Short('k').OverrideDefaultFromEnvar("PRIVATE_KEY_PATH").Required().String()
	p := app.Flag("port", "Listening port").Short('p').Default("80").Int()
	_, err = app.Parse(args)
	if err != nil {
		return
	}

	o = Options{*i, *k, *p}
	return
}
