package publisher_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/minodisk/presigner/options"
	"github.com/minodisk/presigner/publisher"
)

var (
	authJSON = os.Getenv("GOOGLE_AUTH_JSON")
	bucket   = os.Getenv("PRESIGNER_BUCKET")
	account  options.Account
)

func TestMain(m *testing.M) {
	account = options.Account{}
	err := json.Unmarshal([]byte(authJSON), &account)
	if err != nil {
		panic(fmt.Sprintf("fail to initialize Account: %v", err))
	}

	code := m.Run()
	os.Exit(code)
}

func TestPublishAndUpload(t *testing.T) {
	for _, c := range []struct {
		name        string
		pub         publisher.Publisher
		params      publisher.Params
		header      http.Header
		data        string
		disposition string
		ext         string
	}{
		{
			name: "sign without param",
			pub: publisher.Publisher{options.Options{
				ServiceAccount: account,
				Bucket:         bucket,
				Duration:       time.Minute,
			}},
			params:      publisher.Params{},
			header:      http.Header{},
			data:        "foo",
			disposition: "",
			ext:         "",
		},
		{
			name: "setup with ObjectPrefix",
			pub: publisher.Publisher{options.Options{
				ServiceAccount: account,
				Bucket:         bucket,
				Duration:       time.Minute,
				ObjectPrefix:   "foo/",
			}},
			params:      publisher.Params{},
			header:      http.Header{},
			data:        "foo",
			disposition: "",
			ext:         "",
		},
		{
			name: "sign with ContentType",
			pub: publisher.Publisher{options.Options{
				ServiceAccount: account,
				Bucket:         bucket,
				Duration:       time.Minute,
			}},
			params: publisher.Params{
				ContentType: "text/plain",
			},
			header: http.Header{
				"Content-Type": []string{"text/plain"},
			},
			data:        "bar",
			disposition: "",
			ext:         "",
		},
		{
			name: "sign with param Filename",
			pub: publisher.Publisher{options.Options{
				ServiceAccount: account,
				Bucket:         bucket,
				Duration:       time.Minute,
			}},
			params: publisher.Params{
				Filename: "baz.txt",
			},
			header:      http.Header{},
			data:        "baz",
			disposition: "",
			ext:         ".txt",
		},
		{
			name: "upload with Content-Disposition",
			pub: publisher.Publisher{options.Options{
				ServiceAccount: account,
				Bucket:         bucket,
				Duration:       time.Minute,
			}},
			params: publisher.Params{},
			header: http.Header{
				"Content-Disposition": []string{"attachment; filename=qux.txt"},
			},
			data:        "qux",
			disposition: "attachment; filename=qux.txt",
			ext:         "",
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {

			res, err := c.pub.Publish(c.params)
			if err != nil {
				t.Fatalf("fail to publish: %v", err)
			}

			t.Run("SignedURL", func(t *testing.T) {
				if !strings.HasPrefix(res.SignedURL, "https://storage.googleapis.com/presigner/") {
					t.Errorf("SignedURL does not have correct prefix: %s", res.SignedURL)
				}

				req, err := http.NewRequest(http.MethodPut, res.SignedURL, bytes.NewBuffer([]byte(c.data)))
				if err != nil {
					t.Fatal(err)
				}
				req.Header = c.header
				cli := &http.Client{}
				resp, err := cli.Do(req)
				if err != nil {
					t.Fatal(err)
				}
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}
				if resp.StatusCode/100 != 2 {
					t.Errorf("fail to upload:\n%s\n%s", resp.Status, body)
				}
			})

			t.Run("FileURL", func(t *testing.T) {
				{
					got := res.FileURL
					want := fmt.Sprintf("https://presigner.storage.googleapis.com/%s", c.pub.Options.ObjectPrefix)
					if !strings.HasPrefix(got, want) {
						t.Errorf("FileURL does not have correct prefix: got %s, want %s", got, want)
					}
				}
				{
					got := path.Ext(res.FileURL)
					want := c.ext
					if got != want {
						t.Errorf("FileURL does not have correct extension: got %s, want %s", got, want)
					}
				}

				resp, err := http.Get(res.FileURL)
				if err != nil {
					t.Fatalf("fail to get: %v", err)
				}
				{
					want := c.header.Get("Content-Type")
					if want == "" {
						got := resp.Header.Get("Content-Type")
						want := "application/octet-stream"
						if got != want {
							t.Errorf("Content-Type: got %s, want %s", got, want)
						}
					} else {
						got := resp.Header.Get("Content-Type")
						if got != want {
							t.Errorf("Content-Type: got %s, want %s", got, want)
						}
					}
				}
				{
					want := c.disposition
					if want != "" {
						got := resp.Header.Get("Content-Disposition")
						if got != want {
							t.Errorf("Content-Disposition: got %s, want %s", got, want)
						}
					}
				}
				fmt.Println(resp.Header)
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("fail to read: %v", err)
				}
				if resp.StatusCode/100 != 2 {
					t.Fatalf("fail to fetch:\n%s\n%s", resp.Status, body)
				}
				got := string(body)
				if got != c.data {
					t.Errorf("body does not match\n got: %s\nwant: %s", got, c.data)
				}
			})
		})
	}
}
