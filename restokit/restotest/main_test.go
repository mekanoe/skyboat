package main

import (
	"net/http"
	"testing"
)

func TestStartup(t *testing.T) {
	main()

	resp, err := http.Get("http://:4665/test")
	if err != nil || resp.StatusCode != 200 {
		t.Error(err)
		t.Fail()
	}
}
