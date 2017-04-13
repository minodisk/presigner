package options

import "strings"

type Buckets []string

func (b *Buckets) String() string {
	return strings.Join(*b, ", ")
}

func (b *Buckets) Set(s string) error {
	*b = append(*b, s)
	return nil
}

func (bs Buckets) Contains(bucket string) bool {
	if len(bs) == 0 {
		return true
	}
	for _, b := range bs {
		if b == bucket {
			return true
		}
	}
	return false
}
