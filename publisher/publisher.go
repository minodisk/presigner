package publisher

import (
	"fmt"
	"time"

	"google.golang.org/cloud/storage"

	"github.com/go-microservices/presigner/option"
	"github.com/satori/go.uuid"
)

type Publisher struct {
	Bucket      string `json:"bucket"`
	ContentType string `json:"content_type"`
	Filename    string `json:"filename"`
}

type URLSet struct {
	SignedURL string `json:"signed_url"`
	FileURL   string `json:"file_url"`
}

func (p Publisher) Publish(o option.Options) (urlSet URLSet, err error) {
	if len(o.PrivateKey) == 0 {
		err = fmt.Errorf("requires private key bytes")
		return
	}
	if !o.Buckets.Contains(p.Bucket) {
		err = fmt.Errorf("bucket '%s' is not allowed", p.Bucket)
		return
	}

	expiration := time.Now().Add(o.Duration)
	key := uuid.NewV4().String()

	opts := &storage.SignedURLOptions{
		GoogleAccessID: o.GoogleAccessID,
		PrivateKey:     o.PrivateKey,
		Method:         "PUT",
		Expires:        expiration,
		ContentType:    p.ContentType,
	}
	if p.Filename != "" {
		opts.Headers = append(opts.Headers, fmt.Sprintf("Content-Disposition:attachment; filename=%s", p.Filename))
	}
	url, err := storage.SignedURL(p.Bucket, key, opts)
	if err != nil {
		return
	}

	urlSet.SignedURL = url
	urlSet.FileURL = fmt.Sprintf("https://storage.googleapis.com/%s/%s", p.Bucket, key)
	return
}
