package options_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/minodisk/presigner/options"
)

func TestUndefinedFlags(t *testing.T) {
	t.Parallel()
	for i, c := range []struct {
		args []string
	}{
		{
			args: []string{"-xxx"},
		},
	} {
		c := c
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			_, err := options.New(c.args)
			if err == nil {
				t.Error("should error with unknown flags")
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
			"default",
			[]string{},
			options.Options{
				"",
				"/secret/google-auth.json",
				options.Buckets{},
				80,
				time.Minute,
			},
		},
		{
			"multi buckets",
			[]string{
				"-bucket", "bucket-a",
				"-bucket", "bucket-b",
			},
			options.Options{
				"",
				"/secret/google-auth.json",
				options.Buckets{
					"bucket-a",
					"bucket-b",
				},
				80,
				time.Minute,
			},
		},
		{
			"complex",
			[]string{
				"-id", "foo",
				"-key", "./bar.json",
				"-bucket", "bucket-a",
				"-port", "8080",
				"-duration", "1h",
			},
			options.Options{
				"foo",
				"./bar.json",
				options.Buckets{
					"bucket-a",
				},
				8080,
				time.Hour,
			},
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			got, err := options.New(c.args)
			if err != nil {
				t.Fatal("shouldn't error with args: %v", c.args)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("\n got: %+v\nwant: %+v", got, c.want)
			}
		})
	}
}
