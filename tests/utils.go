package tests

import (
	"ChatDanBackend/bootstrap"
	. "ChatDanBackend/models"
	"ChatDanBackend/utils"
	"bytes"
	"github.com/goccy/go-json"
	"github.com/hetiansu5/urlquery"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

var App = bootstrap.InitFiberApp()

type tester struct {
	Token string
}

var (
	defaultTester tester
	userTester    tester
	adminTester   tester
	otherTester   map[int]tester
)

func (tester *tester) testCommon(t *testing.T, method string, route string, statusCode int, isQuery bool, data Map, model any) {
	var requestData []byte
	var err error

	if data != nil {
		if isQuery {
			queryData, err := urlquery.Marshal(data)
			assert.Nilf(t, err, "encode request query")
			route += "?" + string(queryData)
		} else {
			requestData, err = json.Marshal(data)
			assert.Nilf(t, err, "encode request body")
		}
	}
	req, err := http.NewRequest(
		method,
		route,
		bytes.NewBuffer(requestData),
	)
	assert.Nilf(t, err, "constructs http request")
	req.Header.Add("Content-Type", "application/json")
	if tester.Token != "" {
		req.Header.Add("Authorization", "Bearer "+tester.Token)
	}

	res, err := App.Test(req, -1)
	assert.Nilf(t, err, "perform request")
	assert.Equalf(t, statusCode, res.StatusCode, "status code")

	responseBody, err := io.ReadAll(res.Body)
	assert.Nilf(t, err, "decode response")

	if res.StatusCode >= 400 {
		utils.Logger.Error(string(responseBody))
	} else {
		if model != nil {
			err = json.Unmarshal(responseBody, model)
			assert.Nilf(t, err, "decode response")
		}
	}
}

func (tester *tester) testCommonQuery(t *testing.T, method string, route string, statusCode int, data Map, model any) {
	tester.testCommon(t, method, route, statusCode, true, data, model)
}

func (tester *tester) testCommonBody(t *testing.T, method string, route string, statusCode int, data Map, model any) {
	tester.testCommon(t, method, route, statusCode, false, data, model)
}

func (tester *tester) testGet(t *testing.T, route string, statusCode int, data Map, model any) {
	tester.testCommonQuery(t, http.MethodGet, route, statusCode, data, model)
}

func (tester *tester) testPost(t *testing.T, route string, statusCode int, data Map, model any) {
	tester.testCommonBody(t, http.MethodPost, route, statusCode, data, model)
}

func (tester *tester) testPut(t *testing.T, route string, statusCode int, data Map, model any) {
	tester.testCommonBody(t, http.MethodPut, route, statusCode, data, model)
}

func (tester *tester) testDelete(t *testing.T, route string, statusCode int, data Map, model any) {
	tester.testCommonQuery(t, http.MethodDelete, route, statusCode, data, model)
}

func (tester *tester) testPatch(t *testing.T, route string, statusCode int, data Map, model any) {
	tester.testCommonBody(t, http.MethodPatch, route, statusCode, data, model)
}
