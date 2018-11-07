package pagerduty

import (
	"testing"
	"net/http"
	"fmt"
)

type MockHTTPClient struct {
	expectedQuery string
}

func (c MockHTTPClient) Do(request *http.Request) (*http.Response, error) {
	var err error

	resp := &http.Response{StatusCode: 200}
	if request.URL.RawQuery != c.expectedQuery {
		err = fmt.Errorf("query string did not match: %s != %s", c.expectedQuery, request.URL.RawQuery)
	}
	return resp, err
}

func newTestClient(expectedQuery string) HTTPClient {
	return &MockHTTPClient{expectedQuery:expectedQuery}
}

func TestWithQuery(t *testing.T) {
	values := []string{"kylie@example.com"}
	client := NewClient("123", WithCustomClient(newTestClient("query=kylie%40example.com")))
	for _, value := range values {
		_, err := client.ListUsers(WithQuery(value))
		if err != nil {
			t.Log("non-nil error from response: " + err.Error())
			t.Fail()
		}
	}
}
