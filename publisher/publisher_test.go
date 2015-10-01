package publisher_test

import (
	"bytes"
	"io/ioutil"
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
	o, err := option.New([]string{})
	if err != nil {
		t.Fatal(err)
	}
	err = o.ReadPrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	content := []byte("foo")
	p := publisher.Publisher{"signing-test", "text/plain"}
	urlSet, err := p.Publish(o)
	if err != nil {
		t.Fatal(err)
	}
	var reqBody bytes.Buffer
	reqBody.Write(content)
	req, err := http.NewRequest("PUT", urlSet.SignedURL, &reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "text/plain")

	// fmt.Println("----------REQUEST-----------")
	// fmt.Println(string(reqBody.Bytes()))
	// fmt.Println("----------------------------")

	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}

	// fmt.Println("----------RESPONSE----------")
	// fmt.Println(r.StatusCode)
	// fmt.Println(string(resBody))
	// fmt.Println("----------------------------")

	if r.StatusCode != 200 {
		t.Fatal(string(resBody))
	}
}
