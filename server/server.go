package server

import (
	"fmt"
	"net/http"

	"github.com/go-microservices/signing/option"
	"github.com/go-microservices/signing/server/routes"
)

func Serve(o option.Options, privateKey []byte) (err error) {
	http.Handle("/", routes.Index{o, privateKey})
	err = http.ListenAndServe(fmt.Sprintf(":%d", o.Port), nil)
	return
}
