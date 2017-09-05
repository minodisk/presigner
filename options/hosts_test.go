package options_test

import (
	"testing"

	"github.com/minodisk/presigner/options"
)

func TestHostsString(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name  string
		hosts options.Hosts
		want  string
	}{
		{
			"general",
			options.Hosts{
				"foo.com",
				"bar.com",
			},
			"foo.com, bar.com",
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			got := c.hosts.String()
			if got != c.want {
				t.Errorf("got: %s, want: %s", got, c.want)
			}
		})
	}
}
