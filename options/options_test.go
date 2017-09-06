package options_test

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/minodisk/presigner/options"
)

const (
	pathToAccount = "options_test.json"
)

func TestMain(m *testing.M) {
	if err := ioutil.WriteFile(pathToAccount, []byte(`{"client_email": "test@example.com", "private_key": "xxxxxxxxxx\nyyyyyyyyyy\nzzzzzzzzzz\n"}`), 0644); err != nil {
		panic(err)
	}

	code := m.Run()
	os.Remove(pathToAccount)
	os.Exit(code)
}

func TestParseError(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name string
		args []string
	}{
		{
			"with undefined flag",
			[]string{
				"-account", pathToAccount,
				"-xxx",
			},
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			got := &options.Options{}
			if err := got.Parse(c.args); err == nil {
				t.Error("should error")
			}
		})
	}
}

func TestFullfillment(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name string
		args []string
		want *options.Options
	}{
		{
			"complex",
			[]string{
				"-account", pathToAccount,
				"-bucket", "bucket-a",
				"-duration", "1h",
				"-host", "foo.com,bar.com",
				"-port", "8080",
			},
			&options.Options{
				Account: options.Account{
					ClientEmail: "test@example.com",
					PrivateKey:  "xxxxxxxxxx\nyyyyyyyyyy\nzzzzzzzzzz\n",
				},
				Bucket:    "bucket-a",
				Duration:  time.Hour,
				Referrers: []string{"foo.com", "bar.com"},
				Port:      8080,
				Verbose:   false,
			},
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			got := &options.Options{}
			if err := got.Parse(c.args); err != nil {
				t.Fatalf("shouldn't error: %v\nwith args: %v", err, c.args)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("\n got: %+v\nwant: %+v", got, c.want)
			}
		})
	}
}
