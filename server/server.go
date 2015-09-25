package server

import (
	"fmt"
	"net/http"

	"github.com/minodisk/presigner/option"
	"github.com/minodisk/presigner/server/routes"
)

func Serve(o option.Options) (err error) {
	http.Handle("/", routes.Index{o})
	err = http.ListenAndServe(fmt.Sprintf(":%d", o.Port), nil)
	return
}
