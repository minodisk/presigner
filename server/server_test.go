package server_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/minodisk/presigner/options"
	"github.com/minodisk/presigner/server"
)

var (
	Server *httptest.Server
	Client *http.Client

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

	Server = httptest.NewServer(server.Index{opts})
	defer Server.Close()
	Client = &http.Client{}

	code := m.Run()
	os.Remove("google-auth.json")
	os.Exit(code)
}

func TestNotAllowedMethods(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		method string
		err    server.Error
	}{
		{
			"PUT",
			server.Error{
				Error: "PUT method is not allowed",
			},
		},
		{
			"DELETE",
			server.Error{
				Error: "DELETE method is not allowed",
			},
		},
	} {
		c := c
		t.Run(c.method, func(t *testing.T) {
			t.Parallel()
			req, err := http.NewRequest(c.method, Server.URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := Client.Do(req)
			if resp.StatusCode != http.StatusMethodNotAllowed {
				t.Errorf("status code got %v, want %v", resp.StatusCode, http.StatusMethodNotAllowed)
			}

			buf, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("fail to read response body: %v", err)
			}

			var e server.Error
			err = json.Unmarshal(buf, &e)
			if err != nil {
				t.Fatalf("fail to unmarshal JSON: %v", err)
			}
			if !reflect.DeepEqual(e, c.err) {
				t.Errorf("\n got: %+v\nwant: %+v", e, c.err)
			}
		})
	}
}
