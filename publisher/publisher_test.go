package publisher_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

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
	opts, err = options.Options{
		Buckets:  options.Buckets{bucket},
		Duration: time.Minute,
	}.FillAccountWithJSON([]byte(authJSON))
	if err != nil {
		panic(fmt.Sprintf("fail to initialize GoogleAuthKey: %v", err))
	}

	code := m.Run()
	os.Exit(code)
}

func TestUpload(t *testing.T) {
	want := "test"
	fmt.Printf("%+v", publisher.Publisher{
		Filename:    "test.txt",
		Bucket:      bucket,
		ContentType: "text/plain",
	})
	res, err := publisher.Publisher{
		Filename:    "test.txt",
		Bucket:      bucket,
		ContentType: "text/plain",
	}.Publish(opts)
	if err != nil {
		t.Fatalf("fail to publish: %v", err)
	}

	t.Run("SignedURL", func(t *testing.T) {
		req, err := http.NewRequest("PUT", res.SignedURL, bytes.NewBuffer([]byte(want)))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "text/plain")
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
			t.Fatalf("fail to upload:\n%s\n%s", resp.Status, body)
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
