package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"google.golang.org/cloud/storage"

	"github.com/jessevdk/go-flags"
	"github.com/satori/go.uuid"
)

var (
	privateKey []byte
)

type Opts struct {
	Port           int    `short:"p" long:"port" description:"Listening port"`
	GoogleAccessID string `short:"i" long:"id" description:"Google access ID"`
	PrivateKey     string `short:"k" long:"key" description:"Private key"`
}

func main() {
	var opts Opts
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "sign"
	parser.Usage = "[OPTIONS]"
	_, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	key, err := ioutil.ReadFile(opts.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	privateKey = key

	log.Printf("%+v", opts)

	Serve(opts)
}

func Serve(opts Opts) (err error) {
	http.Handle("/", Index{opts})
	err = http.ListenAndServe(fmt.Sprintf(":%d", opts.Port), nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}

type Index struct {
	opts Opts
}

type Req struct {
	Bucket      string `json:"bucket"`
	ContentType string `json:"content_type"`
	MD5         string `json:"md5"`
}

type Resp struct {
	PutURL string `json:"put_url"`
	GetURL string `json:"get_url"`
}

func (i Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	if r.Method == "OPTIONS" {
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		responseError(w, []error{err})
		return
	}

	log.Println("Request body:", string(buf))

	var req Req
	err = json.Unmarshal(buf, &req)
	if err != nil {
		responseError(w, []error{err})
		return
	}

	log.Printf("%+v", req)

	key := uuid.NewV4().String()
	opts := &storage.SignedURLOptions{
		GoogleAccessID: i.opts.GoogleAccessID,
		PrivateKey:     privateKey,
		Method:         "PUT",
		Expires:        time.Now().Add(time.Minute),
		ContentType:    req.ContentType,
		// MD5:            []byte(req.MD5),
	}

	url, err := storage.SignedURL(req.Bucket, key, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("PUT URL : %v\n", url)

	resp := Resp{
		url,
		fmt.Sprintf("https://storage.googleapis.com/%s/%s", req.Bucket, key),
	}
	response(w, resp)
}

func responseError(w http.ResponseWriter, errs []error) {
	var errors []string
	for _, err := range errs {
		log.Println("Error:", err)
		errors = append(errors, err.Error())
	}
	response(w, map[string][]string{"errors": errors})
}

func response(w http.ResponseWriter, body interface{}) {
	buf, err := json.Marshal(body)
	if err != nil {
		log.Fatal("fail to marshal JSON:", err)
		return
	}

	log.Println("Response:", string(buf))

	i, err := w.Write(buf)
	if err != nil {
		log.Fatal("fail to write response:", err)
		return
	}
	log.Printf("write %d bytes", i)
}
