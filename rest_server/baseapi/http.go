package baseapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/google/go-querystring/query"
)

func MakeHttpClient(uri string, auth string, method string, body *bytes.Buffer, queryStr string) (*http.Client, *http.Request) {
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, nil
	}
	if len(auth) > 0 {
		req.Header.Add("Authorization", auth)
	}
	req.Header.Add("Content-Type", "application/json")
	if len(queryStr) > 0 {
		req.URL.RawQuery = queryStr
	}

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}

	return client, req
}

func HttpCall(uri string, auth string, method string, body *bytes.Buffer, queryStruct interface{}) (*base.BaseResponse, error) {
	var v url.Values
	var queryStr string
	if queryStruct != nil {
		v, _ = query.Values(queryStruct)
		queryStr = v.Encode()
	}

	client, req := MakeHttpClient(uri, auth, method, body, queryStr)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
	}

	return ParseResponse(resp)
}

func ParseResponse(resp *http.Response) (*base.BaseResponse, error) {
	decoder := json.NewDecoder(resp.Body)
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		ErrorResp := new(base.BaseResponse)
		err := decoder.Decode(ErrorResp)
		if err != nil {
			return nil, errors.New(resp.Status)
		}
		return nil, err
	}
	baseResp := new(base.BaseResponse)
	err := decoder.Decode(&baseResp)
	if err != nil {
		return nil, err
	}
	return baseResp, err
}
