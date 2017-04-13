package publisher

import (
	"fmt"
	"io/ioutil"
	"time"

	"cloud.google.com/go/storage"

	"github.com/minodisk/presigner/options"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type Publisher struct {
	Filename    string   `json:"filename"`
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

	privateKey, err := ioutil.ReadFile(o.PrivateKeyPath)
	if err != nil {
		return res, errors.Wrap(err, "fail to read private key")
	}
	if !o.Buckets.Contains(p.Bucket) {
		err = fmt.Errorf("the bucket %s is not allowed to sign", p.Bucket)
		return res, err
	}

	expiration := time.Now().Add(o.Duration)
	opts := storage.SignedURLOptions{
		GoogleAccessID: o.GoogleAccessID,
		PrivateKey:     privateKey,
		Method:         "PUT",
		Expires:        expiration,
		ContentType:    p.ContentType,
		Headers:        p.Headers,
		// Headers: append(
		// 	p.Headers,
		// 	fmt.Sprintf("Content-Disposition:attachment; filename=%s", p.Filename),
		// ),
	}
	if p.MD5 != "" {
		opts.MD5 = []byte(p.MD5)
	}
	fmt.Println("MD5:", opts.MD5)

	key := uuid.NewV4().String()
	url, err := storage.SignedURL(p.Bucket, key, &opts)
	if err != nil {
		return res, errors.Wrap(err, "fail to sign")
	}
	res.SignedURL = url
	res.FileURL = fmt.Sprintf("https://storage.googleapis.com/%s/%s", p.Bucket, key)
	return res, nil
}
