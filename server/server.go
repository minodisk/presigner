package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/minodisk/presigner/options"
	"github.com/minodisk/presigner/publisher"
)

func Start(o *options.Options) (err error) {
	if o.Verbose {
		fmt.Printf("Options: %+v\n", o)
	}
	http.Handle("/", Index{o})
	fmt.Printf("listening on port %d\n", o.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", o.Port), nil)
}

type Index struct {
	Options *options.Options
}

func (i Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := func(host, method string, body io.ReadCloser) (*Resp, error) {
		switch method {
		default:
			return nil, NewMethodNotAllowed(method)
		case http.MethodPost:
			b, err := ioutil.ReadAll(body)
			if err != nil {
				return nil, NewBadRequest(err)
			}
			var params publisher.Params
			if err = json.Unmarshal(b, &params); err != nil {
				return nil, NewBadRequest(err)
			}
			if i.Options.Verbose {
				fmt.Printf("Publisher: %+v\n", params)
			}
			pub := publisher.Publisher{i.Options}
			res, err := pub.Publish(params)
			if err != nil {
				return nil, NewBadRequest(err)
			}
			return NewResp(http.StatusOK, res), nil
		}
	}(r.URL.Hostname(), r.Method, r.Body)
	if err != nil {
		if coder, ok := err.(Coder); ok {
			resp = NewErrorResp(coder.Code(), err)
		} else {
			resp = NewErrorResp(500, err)
		}
	}

	if resp.Body == nil {
		w.WriteHeader(resp.Code())
		return
	}

	header := w.Header()
	// header.Set("Access-Control-Allow-Origin", "*")
	// header.Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	// header.Set("Access-Control-Allow-Headers", "Origin, Content-Type")
	header.Set("Content-Type", "application/json")

	b, err := json.Marshal(resp.Body)
	if err != nil {
		log.Printf("fail to marshal JSON: %+v", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(resp.Code())

	if _, err = w.Write(b); err != nil {
		log.Printf("fail to write response body: %+v", err)
		return
	}

	log.Printf("response: %v", string(b))
}
