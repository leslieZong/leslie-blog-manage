package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type APIResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func DecodeAPIResponse(
	t *testing.T,
	recorder *httptest.ResponseRecorder,
) APIResponse {

	t.Helper()

	var response APIResponse

	if err := json.Unmarshal(
		recorder.Body.Bytes(),
		&response,
	); err != nil {

		t.Fatalf(
			"decode API response failed: %v, body=%s",
			err,
			recorder.Body.String(),
		)
	}

	return response
}

func DoJSONRequest(
	t *testing.T,
	router http.Handler,
	method string,
	path string,
	body any,
	token string,
) *httptest.ResponseRecorder {

	t.Helper()

	var bodyReader *bytes.Reader

	if body == nil {
		bodyReader = bytes.NewReader(nil)
	} else {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf(
				"marshal request body failed: %v",
				err,
			)
		}

		bodyReader = bytes.NewReader(data)
	}

	req := httptest.NewRequest(
		method,
		path,
		bodyReader,
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	if token != "" {
		req.Header.Set(
			"Authorization",
			"Bearer "+token,
		)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		req,
	)

	return recorder
}

func AssertStatus(
	t *testing.T,
	recorder *httptest.ResponseRecorder,
	expected int,
) {
	t.Helper()

	if recorder.Code != expected {
		t.Fatalf(
			"expected status %d, got %d, body=%s",
			expected,
			recorder.Code,
			recorder.Body.String(),
		)
	}
}
