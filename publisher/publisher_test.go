package publisher_test

const (
	Bucket = "presigner-test"
)

//
// func TestMain(m *testing.M) {
// 	code := m.Run()
// 	os.Exit(code)
// }
//
// func TestUpload(t *testing.T) {
// 	o, err := option.New([]string{})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = o.ReadPrivateKey()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	content := []byte("foo")
// 	p := publisher.Publisher{"presigner-test", "text/plain"}
// 	urlSet, err := p.Publish(o)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	var reqBody bytes.Buffer
// 	reqBody.Write(content)
// 	req, err := http.NewRequest("PUT", urlSet.SignedURL, &reqBody)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req.Header.Set("Content-Type", "text/plain")
//
// 	// fmt.Println("----------REQUEST-----------")
// 	// fmt.Println(string(reqBody.Bytes()))
// 	// fmt.Println("----------------------------")
//
// 	client := &http.Client{}
// 	r, err := client.Do(req)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	resBody, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	// fmt.Println("----------RESPONSE----------")
// 	// fmt.Println(r.StatusCode)
// 	// fmt.Println(string(resBody))
// 	// fmt.Println("----------------------------")
//
// 	if r.StatusCode != 200 {
// 		t.Fatal(string(resBody))
// 	}
// }
