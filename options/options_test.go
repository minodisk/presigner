package options_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/minodisk/presigner/options"
)

func TestParseError(t *testing.T) {
	t.Parallel()
	for i, c := range []struct {
		name string
		args []string
	}{
		{
			"without -account",
			[]string{},
		},
		{
			"with undefined flag",
			[]string{
				"-account", `{"client_email": "test@example.com", "private_key": "xxxxxxxxxx\nyyyyyyyyyy\nzzzzzzzzzz\n"}`,
				"-xxx",
			},
		},
	} {
		c := c
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			_, err := options.Parse(c.args)
			if err == nil {
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
		want options.Options
	}{
		{
			"with -account",
			[]string{
				"-account", `{"client_email": "test@example.com", "private_key": "xxxxxxxxxx\nyyyyyyyyyy\nzzzzzzzzzz\n"}`,
			},
			options.Options{
				options.Account{
					ClientEmail: "test@example.com",
					PrivateKey:  "xxxxxxxxxx\nyyyyyyyyyy\nzzzzzzzzzz\n",
				},
				options.Buckets{},
				time.Minute,
				80,
				false,
			},
		},
		{
			"multi buckets",
			[]string{
				"-account", `{"client_email": "test@example.com", "private_key": "xxxxxxxxxx\nyyyyyyyyyy\nzzzzzzzzzz\n"}`,
				"-bucket", "bucket-a",
				"-bucket", "bucket-b",
			},
			options.Options{
				options.Account{
					ClientEmail: "test@example.com",
					PrivateKey:  "xxxxxxxxxx\nyyyyyyyyyy\nzzzzzzzzzz\n",
				},
				options.Buckets{
					"bucket-a",
					"bucket-b",
				},
				time.Minute,
				80,
				false,
			},
		},
		{
			"complex",
			[]string{
				"-account", `{"client_email": "test@example.com", "private_key": "xxxxxxxxxx\nyyyyyyyyyy\nzzzzzzzzzz\n"}`,
				"-bucket", "bucket-a",
				"-duration", "1h",
				"-port", "8080",
			},
			options.Options{
				options.Account{
					ClientEmail: "test@example.com",
					PrivateKey:  "xxxxxxxxxx\nyyyyyyyyyy\nzzzzzzzzzz\n",
				},
				options.Buckets{
					"bucket-a",
				},
				time.Hour,
				8080,
				false,
			},
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			got, err := options.Parse(c.args)
			if err != nil {
				t.Fatalf("shouldn't error: %v\nwith args: %v", err, c.args)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("\n got: %+v\nwant: %+v", got, c.want)
			}
		})
	}
}
