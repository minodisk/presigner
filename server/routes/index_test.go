package routes_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-microservices/signing/option"
	"github.com/go-microservices/signing/server/routes"
)

var (
	server *httptest.Server
	client *http.Client
)

func TestMain(m *testing.M) {
	server = httptest.NewServer(routes.Index{option.Options{}})
	client = &http.Client{}
	code := m.Run()
	server.Close()
	os.Exit(code)
}

func TestNotAllowedMethods(t *testing.T) {
	for _, method := range []string{"GET", "PUT", "DELETE"} {
		req, err := http.NewRequest(method, server.URL, nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := client.Do(req)
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("method %s should be rejected", method)
		}

		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		var body struct {
			Errors []string
		}
		err = json.Unmarshal(buf, &body)
		if err != nil {
			t.Fatal(err)
		}

		if len(body.Errors) != 1 {
			t.Errorf("errors should be 1")
		} else {
			if body.Errors[0] != "POST method is allowed" {
				t.Errorf("error is not valid")
			}
		}
	}
}
