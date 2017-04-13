package options_test

import (
	"reflect"
	"testing"

	"github.com/minodisk/presigner/options"
)

func TestString(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name    string
		buckets options.Buckets
		want    string
	}{
		{
			"empty",
			options.Buckets{},
			"",
		},
		{
			"single",
			options.Buckets{
				"foo",
			},
			"foo",
		},
		{
			"multi",
			options.Buckets{
				"foo",
				"bar",
			},
			"foo, bar",
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			got := c.buckets.String()
			if got != c.want {
				t.Errorf("\n got: %v\nwant: %v", got, c.want)
			}
		})
	}
}

func TestSet(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name   string
		inputs []string
		want   *options.Buckets
	}{
		{
			"not set",
			[]string{},
			&options.Buckets{},
		},
		{
			"single",
			[]string{
				"foo",
			},
			&options.Buckets{
				"foo",
			},
		},
		{
			"multi",
			[]string{
				"foo",
				"bar",
				"baz",
			},
			&options.Buckets{
				"foo",
				"bar",
				"baz",
			},
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			got := &options.Buckets{}
			for _, input := range c.inputs {
				got.Set(input)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("\n got: %+v\nwant: %+v", got, c.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name    string
		buckets options.Buckets
		bucket  string
		want    bool
	}{
		{
			"empty buckets allows any bucket",
			options.Buckets{},
			"foo",
			true,
		},
		{
			"single",
			options.Buckets{
				"foo",
			},
			"foo",
			true,
		},
		{
			"single",
			options.Buckets{
				"foo",
			},
			"bar",
			false,
		},
		{
			"multi",
			options.Buckets{
				"foo",
				"bar",
			},
			"foo",
			true,
		},
		{
			"multi",
			options.Buckets{
				"foo",
				"bar",
			},
			"bar",
			true,
		},
		{
			"multi",
			options.Buckets{
				"foo",
				"bar",
			},
			"baz",
			false,
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			got := c.buckets.Contains(c.bucket)
			if got != c.want {
				t.Errorf("\n got: %v\nwant: %v", got, c.want)
			}
		})
	}
}
