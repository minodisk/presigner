package publisher

import (
	"fmt"
	"time"

	"github.com/go-microservices/signing-gcs/option"
	"github.com/satori/go.uuid"
)

type Req struct {
	ContentType string
	Size        int
}

type Resp struct {
	URL    string            `json:"url"`
	Fields map[string]string `json:"fields"`
	Errors []string          `json:"errors"`
}

func Publish(options option.Options, req Req) (resp Resp, err error) {
	key := uuid.NewV4().String()

	s, err := NewSign(options.SecretAccessKey, options.Bucket, key, req.Size, req.Size, time.Now().Add(time.Duration(options.Duration)))
	if err != nil {
		return
	}

	resp = Resp{
		URL: fmt.Sprintf("https://%s.s3.amazonaws.com/", options.Bucket),
		Fields: map[string]string{
			"AWSAccessKeyId": options.AccessKeyID,
			"policy":         s.Policy,
			"signature":      s.Signature,
			"key":            key,
		},
	}

	return
}
