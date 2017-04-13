package server_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
)

func TestMain(m *testing.M) {
	o, err := options.New([]string{})
	if err != nil {
		log.Fatal(err)
	}

	Server = httptest.NewServer(server.Index{o})
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
