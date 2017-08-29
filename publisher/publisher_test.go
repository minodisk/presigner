package publisher_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/minodisk/presigner/options"
	"github.com/minodisk/presigner/publisher"
)

var (
	authJSON = os.Getenv("GOOGLE_AUTH_JSON")
	bucket   = os.Getenv("PRESIGNER_BUCKET")
	opts     options.Options
)

func TestMain(m *testing.M) {
	var err error
	if err := ioutil.WriteFile("google-auth.json", []byte(authJSON), 0644); err != nil {
		panic(err)
	}
	opts, err = options.Parse([]string{
		"-account", "google-auth.json",
		"-bucket", bucket,
	})
	if err != nil {
		panic(fmt.Sprintf("fail to initialize Account: %v", err))
	}

	code := m.Run()
	os.Remove("google-auth.json")
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
			name:        "sign without param",
			pub:         publisher.Publisher{opts},
			params:      publisher.Params{},
			header:      http.Header{},
			data:        "foo",
			disposition: "",
			ext:         "",
		},
		{
			name: "sign with ContentType",
			pub:  publisher.Publisher{opts},
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
			pub:  publisher.Publisher{opts},
			params: publisher.Params{
				Filename: "baz.txt",
			},
			header:      http.Header{},
			data:        "baz",
			disposition: "",
			ext:         ".txt",
		},
		{
			name:   "upload with Content-Disposition",
			pub:    publisher.Publisher{opts},
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
				if !strings.HasPrefix(res.FileURL, "https://presigner.storage.googleapis.com/") {
					t.Errorf("FileURL does not have correct prefix: %s", res.FileURL)
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
