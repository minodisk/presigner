package publisher

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"

	"github.com/minodisk/presigner/options"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type Params struct {
	ContentType string   `json:"content_type"`
	Filename    string   `json:"filename"`
	Headers     []string `json:"headers"`
	MD5         string   `json:"md5"`
}

type Result struct {
	SignedURL string `json:"signed_url"`
	FileURL   string `json:"file_url"`
}

type Publisher struct {
	Options options.Options
}

func (p Publisher) Publish(params Params) (Result, error) {
	var res Result

	expiration := time.Now().Add(p.Options.Duration)
	opts := storage.SignedURLOptions{
		GoogleAccessID: p.Options.ServiceAccount.ClientEmail,
		PrivateKey:     []byte(p.Options.ServiceAccount.PrivateKey),
		Method:         http.MethodPut,
		Expires:        expiration,
		ContentType:    params.ContentType,
		Headers:        params.Headers,
	}
	if params.MD5 != "" {
		opts.MD5 = []byte(params.MD5)
	}

	key := p.Options.ObjectPrefix + uuid.NewV4().String() + filepath.Ext(params.Filename)
	if p.Options.Verbose {
		fmt.Printf("Sign with:\n  Key: %s\n  SingedURLOptions: %+v\n", key, opts)
	}
	signed, err := storage.SignedURL(p.Options.Bucket, key, &opts)
	if err != nil {
		return res, errors.Wrap(err, "fail to sign")
	}

	res.SignedURL = signed
	res.FileURL = fmt.Sprintf("https://%s.storage.googleapis.com/%s", p.Options.Bucket, key)
	if p.Options.Verbose {
		fmt.Printf("Result: %+v\n", res)
	}
	return res, nil
}
