package publisher_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
		"-verbose",
	})
	if err != nil {
		panic(fmt.Sprintf("fail to initialize Account: %v", err))
	}

	code := m.Run()
	os.Remove("google-auth.json")
	os.Exit(code)
}

func TestUpload(t *testing.T) {
	want := "test"
	pub := publisher.Publisher{
		Bucket:      bucket,
		ContentType: "text/plain",
	}
	res, err := pub.Publish(opts)
	if err != nil {
		t.Fatalf("fail to publish: %v", err)
	}

	t.Run("SignedURL", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPut, res.SignedURL, bytes.NewBuffer([]byte(want)))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "text/plain")
		req.Header.Set("Content-Disposition", "attachment; filename=\"test.txt\"")
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
		resp, err := http.Get(res.FileURL)
		if err != nil {
			t.Fatalf("fail to get: %v", err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("fail to read: %v", err)
		}
		if resp.StatusCode/100 != 2 {
			t.Fatalf("fail to fetch:\n%s\n%s", resp.Status, body)
		}
		got := string(body)
		if got != want {
			t.Errorf("\n got: %s\nwant: %s", got, want)
		}
	})
}
