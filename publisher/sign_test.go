package publisher_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-microservices/signing/publisher"
)

func TestNewSign(t *testing.T) {
	fmt.Printf("%+v", time.Unix(0, 0))
	s, err := publisher.NewSign("X", "Y", "Z", 100, 100, time.Unix(0, 0))
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v", s)

	if s.Policy != "eyJjb25kaXRpb25zIjpbeyJidWNrZXQiOiJZIn0seyJrZXkiOiJaIn0sWyJjb250ZW50LWxlbmd0aC1yYW5nZSIsMTAwLDEwMF1dLCJleHBpcmF0aW9uIjoiMTk3MC0wMS0wMVQwOTowMDowMCswOTowMCJ9" {
		t.Error("wrong policy")
	}
	if s.Signature != "AwJdS1U0z8JIitZwWCtSRjzU0eA=" {
		t.Error("wrong signature")
	}
}
