package option_test

import (
	"testing"

	"github.com/go-microservices/presigner/option"
)

// func TestRequired(t *testing.T) {
// 	_, err := option.New([]string{})
// 	if err == nil {
// 		t.Error("should error without required flags")
// 	}
// }

func TestUnknown(t *testing.T) {
	_, err := option.New([]string{"-x"})
	if err == nil {
		t.Error("should error with unknown flags")
	}
}

func TestDefault(t *testing.T) {
	o, err := option.New([]string{
		"-i", "AAAA",
		"-k", "BBBB",
	})
	if err != nil {
		t.Fatal(err)
	}

	if o.GoogleAccessID != "AAAA" {
		t.Error("wrong GoogleAccessID:", o.GoogleAccessID)
	}
	if o.PrivateKeyPath != "BBBB" {
		t.Error("wrong PrivateKeyPath:", o.PrivateKeyPath)
	}
	if o.Port != 80 {
		t.Error("wrong Port:", o.Port)
	}
}
