package server

import (
	"fmt"
	"net/http"

	"github.com/go-microservices/signing/option"
	"github.com/go-microservices/signing/server/routes"
)

func Serve(options option.Options) (err error) {
	http.Handle("/", routes.Index{options})
	err = http.ListenAndServe(fmt.Sprintf(":%d", options.Port), nil)
	if err != nil {
		return
	}
	return
}
