package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/minodisk/presigner/options"
	"github.com/minodisk/presigner/publisher"
	"github.com/minodisk/presigner/server"
)

var (
	Server *httptest.Server
	Client *http.Client

	authJSON = os.Getenv("GOOGLE_AUTH_JSON")
	bucket   = os.Getenv("PRESIGNER_BUCKET")
)

func TestMain(m *testing.M) {
	account := options.Account{}
	err := json.Unmarshal([]byte(authJSON), &account)
	if err != nil {
		panic(fmt.Sprintf("fail to initialize Account: %v", err))
	}
	fmt.Println(account)

	Server = httptest.NewServer(server.Index{&options.Options{
		Account:  account,
		Bucket:   bucket,
		Duration: time.Minute,
		Port:     80,
		Verbose:  true,
	}})
	defer Server.Close()
	Client = &http.Client{}

	code := m.Run()
	os.Exit(code)
}

func TestNotAllowedMethods(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		method string
		err    server.Error
	}{
		{
			http.MethodGet,
			server.Error{
				Error: "GET method is not allowed",
			},
		},
		{
			http.MethodPut,
			server.Error{
				Error: "PUT method is not allowed",
			},
		},
		{
			http.MethodPatch,
			server.Error{
				Error: "PATCH method is not allowed",
			},
		},
		{
			http.MethodDelete,
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

func TestPOST(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name   string
		params publisher.Params
	}{
		{
			name:   "without param",
			params: publisher.Params{},
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			b, err := json.Marshal(c.params)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, Server.URL, bytes.NewBuffer(b))
			if err != nil {
				t.Fatal(err)
			}

			resp, err := Client.Do(req)

			buf, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("fail to read response body: %v", err)
			}

			if resp.StatusCode != http.StatusOK {
				t.Errorf("status code got %v, want %v: %s", resp.StatusCode, http.StatusOK, buf)
			}

			var res publisher.Result
			err = json.Unmarshal(buf, &res)
			if err != nil {
				t.Fatalf("fail to unmarshal JSON: %v", err)
			}
			if res.SignedURL == "" || res.FileURL == "" {
				t.Errorf("result doesn't filled: %+v", res)
			}
		})
	}
}
