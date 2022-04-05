package api_inno_log

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
)

type api_kind int

const (
	Api_post_account_auth   = 0 // web 사용자 로그 : PostAccountAuth
	Api_post_member_auth    = 1 // app 사용자 로그 : PostMemberAuth
	Api_post_account_coins  = 2 // mod account coins 로그 : PostAccountCoins
	Api_post_member_points  = 3 // mod member point 로그 : PostMemberPoints
	Api_post_exchange_goods = 4 // exchange goods 로그 : PostExchangeGoods
)

type ApiInfo struct {
	ApiType          api_kind
	Uri              string
	Method           string
	ResponseFuncType func() interface{}
	client           *http.Client
}

var ApiList = map[api_kind]ApiInfo{
	Api_post_account_auth:   ApiInfo{ApiType: Api_post_account_auth, Uri: "/account/auth", Method: "POST", ResponseFuncType: func() interface{} { return new(Common) }, client: NewClient()},
	Api_post_member_auth:    ApiInfo{ApiType: Api_post_member_auth, Uri: "/member/auth", Method: "POST", ResponseFuncType: func() interface{} { return new(Common) }, client: NewClient()},
	Api_post_account_coins:  ApiInfo{ApiType: Api_post_account_coins, Uri: "/account/coins", Method: "POST", ResponseFuncType: func() interface{} { return new(Common) }, client: NewClient()},
	Api_post_member_points:  ApiInfo{ApiType: Api_post_member_points, Uri: "/member/points", Method: "POST", ResponseFuncType: func() interface{} { return new(Common) }, client: NewClient()},
	Api_post_exchange_goods: ApiInfo{ApiType: Api_post_exchange_goods, Uri: "/account/exchangegoods", Method: "POST", ResponseFuncType: func() interface{} { return new(Common) }, client: NewClient()},
}

func NewClient() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxIdleConnsPerHost = 100
	t.IdleConnTimeout = 30 * time.Second
	t.DisableKeepAlives = false
	t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}
	return client
}

func MakeHttp(callUrl string, auth string, method string, body *bytes.Buffer, queryStr string) *http.Request {
	req, err := http.NewRequest(method, callUrl, body)
	if err != nil {
		return nil
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	if len(auth) > 0 {
		req.Header.Add("Authorization", "Bearer "+auth)
	}
	if len(queryStr) > 0 {
		req.URL.RawQuery = queryStr
	}
	return req
}

func HttpCall(client *http.Client, callUrl string, auth string, method string, kind api_kind, body *bytes.Buffer, queryStruct interface{}, response interface{}) (interface{}, error) {
	var v url.Values
	var queryStr string
	if queryStruct != nil {
		v, _ = query.Values(queryStruct)
		queryStr = v.Encode()
	}

	req := MakeHttp(callUrl, auth, method, body, queryStr)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := ParseResponse(resp, kind, response)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ParseResponse(resp *http.Response, kind api_kind, response interface{}) (interface{}, error) {
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, errors.New(resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(response)
	if err != nil {
		return nil, err
	}
	return response, err
}
