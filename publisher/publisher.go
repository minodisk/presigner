package publisher

import (
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/storage"

	"github.com/minodisk/presigner/options"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type Publisher struct {
	Bucket      string   `json:"bucket"`
	ContentType string   `json:"content_type"`
	Headers     []string `json:"headers"`
	MD5         string   `json:"md5"`
}

type Result struct {
	SignedURL string `json:"signed_url"`
	FileURL   string `json:"file_url"`
}

func (p Publisher) Publish(o options.Options) (Result, error) {
	var res Result

	if p.Bucket == "" {
		return res, fmt.Errorf("bucket is empty")
	}
	if !o.Buckets.Contains(p.Bucket) {
		return res, fmt.Errorf("the bucket %s is not allowed to sign", p.Bucket)
	}

	expiration := time.Now().Add(o.Duration)
	opts := storage.SignedURLOptions{
		GoogleAccessID: o.ServiceAccount.ClientEmail,
		PrivateKey:     []byte(o.ServiceAccount.PrivateKey),
		Method:         http.MethodPut,
		Expires:        expiration,
		ContentType:    p.ContentType,
		Headers:        p.Headers,
	}
	if p.MD5 != "" {
		opts.MD5 = []byte(p.MD5)
	}

	key := uuid.NewV4().String()
	if o.Verbose {
		fmt.Printf("Sign with:\n  Bucket: %s\n  Key: %s\n  SingedURLOptions: %+v\n", p.Bucket, key, opts)
	}
	signed, err := storage.SignedURL(p.Bucket, key, &opts)
	if err != nil {
		return res, errors.Wrap(err, "fail to sign")
	}

	res.SignedURL = signed
	res.FileURL = fmt.Sprintf("https://%s.storage.googleapis.com/%s", p.Bucket, key)
	if o.Verbose {
		fmt.Printf("Result: %+v\n", res)
	}
	return res, nil
}
