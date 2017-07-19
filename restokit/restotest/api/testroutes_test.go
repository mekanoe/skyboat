package api

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"skyboat.io/x/restokit"
)

func testGetVersion(url, version string, client *http.Client, t *testing.T) (*http.Response, error) {
	req, err := http.NewRequest("GET", "http://localhost:40000"+url, bytes.NewBuffer([]byte("")))
	if err != nil {
		t.Errorf("testGetVersion - %s %s\n  %v", url, version, err)
		t.Fail()
		return nil, err
	}

	req.Header.Add("Accept", fmt.Sprintf("application/vnd.spln.v%s+json", version))

	rsp, err := client.Do(req)
	return rsp, err
}

func TestVersionedRoutes(t *testing.T) {
	resto, client := restokit.ScaffoldHTTP()
	FetchAPIRoutes(resto.Router)
	go resto.Start()

	// default
	rsp, err := testGetVersion("/test", "1", client, t)
	if err != nil {
		t.Error(err)
		t.FailNow()
		return
	}

	if rsp.StatusCode != 200 {
		t.Errorf("status code not 200, was %d", rsp.StatusCode)
		t.FailNow()
		return
	}

	if rsp.Header.Get("SPLN-API-Version") != "v1" {
		t.Errorf("expected API Version header to be v1, got: %s", rsp.Header.Get("SPLN-API-Version"))
		t.Fail()
	}

	rsp, err = testGetVersion("/test", "2", client, t)
	if err != nil {
		t.Error(err)
		t.FailNow()
		return
	}

	if rsp.StatusCode != 200 {
		t.Errorf("status code not 200, was %d", rsp.StatusCode)
		t.FailNow()
		return
	}

	if rsp.Header.Get("SPLN-API-Version") != "v2" {
		t.Errorf("expected API Version header to be v2, got: %s", rsp.Header.Get("SPLN-API-Version"))
		t.Fail()
	}

	restokit.TeardownHTTP(resto)
}
