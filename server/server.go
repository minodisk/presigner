package server

import (
	"fmt"
	"net/http"

	"github.com/go-microservices/signing/option"
	"github.com/go-microservices/signing/server/routes"
)

func Serve(o option.Options) (err error) {
	http.Handle("/", routes.Index{o})
	err = http.ListenAndServe(fmt.Sprintf(":%d", o.Port), nil)
	return
}
