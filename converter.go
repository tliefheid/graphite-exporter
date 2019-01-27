package main

import (
	"io/ioutil"
	"net/http"
)

func getString(body *http.Response) (s string) {
	return string(getBytes(body))
}

func getBytes(body *http.Response) (bytes []byte) {
	b, _ := ioutil.ReadAll(body.Body)
	return b
}
