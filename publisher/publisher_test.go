package publisher_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/go-microservices/signing/option"
	"github.com/go-microservices/signing/publisher"
)

const (
	Bucket = "signing-test"
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

	content := []byte("foo")
	p := publisher.Publisher{"text/plain", len(content)}
	form, err := p.Publish(o.AccessKeyID, o.SecretAccessKey, o.Bucket, o.Duration)
	if err != nil {
		t.Fatal(err)
	}
	var reqBody bytes.Buffer
	w := multipart.NewWriter(&reqBody)
	for fieldname, value := range form.Fields {
		err := w.WriteField(fieldname, value)
		if err != nil {
			t.Fatal(err)
		}
	}
	part, err := w.CreateFormFile("file", "foo.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, err = part.Write(content)
	if err != nil {
		t.Fatal(err)
	}
	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", form.URL, &reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	// req.Header.Set("Content-Length", fmt.Sprintf("%d", reqBody.Len()))

	fmt.Println("----------REQUEST-----------")
	fmt.Println(string(reqBody.Bytes()))
	fmt.Println("----------------------------")

	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("----------RESPONSE----------")
	fmt.Println(r.StatusCode)
	fmt.Println(string(resBody))
	fmt.Println("----------------------------")

	if r.StatusCode != 200 {
		t.Fatal(string(resBody))
	}
}
