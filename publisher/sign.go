package publisher

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

type Sign struct {
	Policy    string
	Signature string
}

func NewSign(secretAccessKey, bucket, key string, min, max int, expiration time.Time) (s Sign, err error) {
	policyDocument := map[string]interface{}{
		"expiration": expiration,
		"conditions": []interface{}{
			map[string]string{"bucket": bucket},
			map[string]string{"key": key},
			[]interface{}{"content-length-range", min, max},
		},
	}
	policyJSON, err := json.Marshal(policyDocument)
	if err != nil {
		return
	}
	s.Policy = strings.Replace(serialize(policyJSON), "\n", "", -1)

	hash := hmac.New(sha1.New, []byte(secretAccessKey))
	hash.Write([]byte(s.Policy))
	mac := hash.Sum(nil)
	s.Signature = strings.Replace(serialize(mac), "\n", "", -1)

	return
}

func serialize(src []byte) string {
	var dest bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &dest)
	encoder.Write(src)
	encoder.Close()
	return string(dest.Bytes())
}
