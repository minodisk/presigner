package publisher_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/minodisk/presigner/options"
	"github.com/minodisk/presigner/publisher"
)

const (
	PrivateKeyPath = "./google-auth.pem"
)

var (
	GoogleAccessID = os.Getenv("GOOGLE_AUTH_EMAIL")
	PrivateKey     = os.Getenv("GOOGLE_AUTH_KEY")
	Bucket         = os.Getenv("BUCKET")
)

func TestMain(m *testing.M) {
	pem := strings.Replace(PrivateKey, `\n`, "\n", -1)
	if err := ioutil.WriteFile(PrivateKeyPath, []byte(pem), 0664); err != nil {
		panic(err)
	}
	code := m.Run()
	if err := os.Remove(PrivateKeyPath); err != nil {
		panic(err)
	}
	os.Exit(code)
}

func TestUpload(t *testing.T) {
	want := "test"
	res, err := publisher.Publisher{
		Filename:    "test.txt",
		Bucket:      Bucket,
		ContentType: "text/plain",
	}.Publish(options.Options{
		GoogleAuthEmail: GoogleAccessID,
		GoogleAuthKey:   PrivateKeyPath,
		Buckets:         options.Buckets{Bucket},
		Duration:        time.Minute,
	})
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
