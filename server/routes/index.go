package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-microservices/signing-gcs/option"
	"github.com/go-microservices/signing-gcs/publisher"
)

type Index struct {
	Options option.Options
}

func (i Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		responseError(w, http.StatusMethodNotAllowed, []error{fmt.Errorf("POST method is allowed")})
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, http.StatusBadRequest, []error{err})
		return
	}

	var req publisher.Req
	err = json.Unmarshal(buf, &req)
	if err != nil {
		responseError(w, http.StatusBadRequest, []error{err})
		return
	}

	resp, err := publisher.Publish(i.Options, req)
	if err != nil {
		responseError(w, http.StatusBadRequest, []error{err})
		return
	}

	response(w, http.StatusOK, resp)
}

func response(w http.ResponseWriter, code int, resp publisher.Resp) {
	j, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("fail to marshal JSON response body: err=%+v, resp=%+v", err, resp)
		return
	}

	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	i, err := w.Write(j)
	if err != nil {
		log.Fatal("fail to write response body: err=%+v, body=%+v", err, j)
	}
	log.Printf("write %d bytes", i)
}

func responseError(w http.ResponseWriter, code int, errs []error) {
	strs := make([]string, len(errs))
	for i, err := range errs {
		strs[i] = err.Error()
	}
	response(w, code, publisher.Resp{Errors: strs})
}
