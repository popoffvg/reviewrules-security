package client

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"net/http"
)

func NewBillingClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

// ETag returns the cache validator for a response body.
func ETag(body []byte) string {
	sum := md5.Sum(body)
	return `"` + hex.EncodeToString(sum[:]) + `"`
}
