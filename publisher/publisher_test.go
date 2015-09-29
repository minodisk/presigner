package publisher_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/go-microservices/signing/option"
	"github.com/go-microservices/signing/publisher"
)

const (
	Bucket = "signing-s3-test"
)

func TestMain(m *testing.M) {
	// Create bucket for test
	// config := aws.NewConfig()
	// config.Region = aws.String(os.Getenv("AWS_REGION"))
	// client := s3.New(config)
	// co, err := client.CreateBucket(&s3.CreateBucketInput{
	// 	Bucket: aws.String(Bucket),
	// 	ACL:    aws.String("public-read-write"),
	// })
	// if err != nil {
	// 	log.Fatalf("fail to create bucket: %s", err)
	// }
	// log.Printf("bucket created: %s", co)

	code := m.Run()

	// Destroy bucket for test
	// do, err := client.DeleteBucket(&s3.DeleteBucketInput{
	// 	Bucket: aws.String(Bucket),
	// })
	// if err != nil {
	// 	log.Fatalf("fail to delete bucket: %s", err)
	// }
	// log.Printf("bucket deleted: %s", do)

	os.Exit(code)
}

func TestUpload(t *testing.T) {
	o, err := option.New([]string{
		"-b", Bucket,
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%+v", o)

	content := []byte("foo")
	resp, err := publisher.Publish(o, publisher.Req{"text/plain", len(content)})
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%+v", resp.Fields)

	// req, err := http.NewRequest("POST", resp.URL, nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	var body bytes.Buffer

	w := multipart.NewWriter(&body)
	part, err := w.CreateFormFile("file", "foo.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, err = part.Write(content)
	if err != nil {
		t.Fatal(err)
	}
	for fieldname, value := range resp.Fields {
		err := w.WriteField(fieldname, value)
		if err != nil {
			t.Fatal(err)
		}
	}
	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}

	// file, header, err := req.FormFile("file")
	// // r, err := http.PostForm(resp.URL, resp.Values)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	req, err := http.NewRequest("POST", resp.URL, &body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", w.FormDataContentType())

	// fmt.Println("RESPONSE+++++")
	// fmt.Println(string(buf))
	// fmt.Println("+++++++++++++")

	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("RESPONSE+++++")
	fmt.Println(string(buf))
	fmt.Println("+++++++++++++")
}
