package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-microservices/presigner/option"
	"github.com/go-microservices/presigner/publisher"
)

type Index struct {
	Options option.Options
}

func (i Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !(r.Method == "OPTIONS" || r.Method == "POST") {
		responseError(w, http.StatusMethodNotAllowed, []error{fmt.Errorf("POST method is allowed")})
		return
	}

	header := w.Header()
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	header.Set("Access-Control-Allow-Headers", "Origin, Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, http.StatusBadRequest, []error{err})
		return
	}

	var p publisher.Publisher
	err = json.Unmarshal(reqBody, &p)
	if err != nil {
		responseError(w, http.StatusBadRequest, []error{err})
		return
	}

	urlSet, err := p.Publish(i.Options)
	if err != nil {
		responseError(w, http.StatusBadRequest, []error{err})
		return
	}

	response(w, http.StatusOK, urlSet)
}

func response(w http.ResponseWriter, code int, body interface{}) {
	b, err := json.Marshal(body)
	if err != nil {
		log.Printf("fail to marshal response body as JSON: err=%+v, body=%+v", err, body)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(b)
	if err != nil {
		log.Printf("fail to write response body: err=%+v, body=%+v", err, b)
		return
	}
	log.Printf("success response:", string(b))
}

func responseError(w http.ResponseWriter, code int, errs []error) {
	strs := make([]string, len(errs))
	for i, err := range errs {
		strs[i] = err.Error()
	}
	response(w, code, map[string][]string{"errors": strs})
}
