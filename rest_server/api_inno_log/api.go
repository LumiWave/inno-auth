package api_inno_log

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (o *ServerInfo) PostAccountAuth(req *AccountAuthLog) (*Common, error) {
	api := ApiList[Api_post_account_auth]
	callUrl := fmt.Sprintf("%s%s%s", o.IntHostUri, o.IntVer, api.Uri)

	pbytes, _ := json.Marshal(req)
	buff := bytes.NewBuffer(pbytes)

	data, err := HttpCall(api.client, callUrl, o.ApiKey, api.Method, api.ApiType, buff, nil, api.ResponseFuncType())
	if err != nil {
		return nil, err
	}

	return data.(*Common), nil
}

func (o *ServerInfo) PostMemberAuth(req *MemberAuthLog) (*Common, error) {
	api := ApiList[Api_post_member_auth]
	callUrl := fmt.Sprintf("%s%s%s", o.IntHostUri, o.IntVer, api.Uri)

	pbytes, _ := json.Marshal(req)
	buff := bytes.NewBuffer(pbytes)

	data, err := HttpCall(api.client, callUrl, o.ApiKey, api.Method, api.ApiType, buff, nil, api.ResponseFuncType())
	if err != nil {
		return nil, err
	}

	return data.(*Common), nil
}

func (o *ServerInfo) PostAccountCoins(req *AccountCoinLog) (*Common, error) {
	api := ApiList[Api_post_account_coins]
	callUrl := fmt.Sprintf("%s%s%s", o.IntHostUri, o.IntVer, api.Uri)

	pbytes, _ := json.Marshal(req)
	buff := bytes.NewBuffer(pbytes)

	data, err := HttpCall(api.client, callUrl, o.ApiKey, api.Method, api.ApiType, buff, nil, api.ResponseFuncType())
	if err != nil {
		return nil, err
	}

	return data.(*Common), nil
}

func (o *ServerInfo) PostMemberPoints(req *MemberPointsLog) (*Common, error) {
	api := ApiList[Api_post_account_coins]
	callUrl := fmt.Sprintf("%s%s%s", o.IntHostUri, o.IntVer, api.Uri)

	pbytes, _ := json.Marshal(req)
	buff := bytes.NewBuffer(pbytes)

	data, err := HttpCall(api.client, callUrl, o.ApiKey, api.Method, api.ApiType, buff, nil, api.ResponseFuncType())
	if err != nil {
		return nil, err
	}

	return data.(*Common), nil
}

func (o *ServerInfo) PostExchangeGoods(req *ExchangeGoodsLog) (*Common, error) {
	api := ApiList[Api_post_exchange_goods]
	callUrl := fmt.Sprintf("%s%s%s", o.IntHostUri, o.IntVer, api.Uri)

	pbytes, _ := json.Marshal(req)
	buff := bytes.NewBuffer(pbytes)

	data, err := HttpCall(api.client, callUrl, o.ApiKey, api.Method, api.ApiType, buff, nil, api.ResponseFuncType())
	if err != nil {
		return nil, err
	}

	return data.(*Common), nil
}
