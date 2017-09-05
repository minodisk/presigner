package options_test

import (
	"testing"

	"github.com/minodisk/presigner/options"
)

func TestAccountString(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name    string
		account options.Account
		want    string
	}{
		{
			"general",
			options.Account{
				ClientEmail: "bar@example.com",
				PrivateKey:  "XXXXXXXXX\nYYYYYYYYYY",
				ProjectID:   "baz",
			},
			"bar@example.com",
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			got := c.account.String()
			if got != c.want {
				t.Errorf("got: %s, want: %s", got, c.want)
			}
		})
	}
}
