package apis

import (
	"fmt"
	"gin-web/core/netlib"
	urlNet "net/url"
	"strings"
	"time"
)

type Api interface {
	Uri() string
	BandData(args ...interface{}) *map[string]interface{}
	GetTimeOut() time.Duration
	SetUri(uri string)
	SetToken(cid int) string
}

//get 请求 args[0]:map[string]string, args[1]:time.Duration
func Get(uri string, args ...interface{}) *netlib.HTTPRequest {
	api := GetApi(uri)
	url := api.Uri()

	data := api.BandData(args...)
	if len(*data) > 0 {
		//参数 url-encode 处理
		queryStr := urlNet.Values{}
		for k, v := range *data {
			queryStr.Add(k, fmt.Sprintf("%v", v))
		}
		if strings.Contains(url, "?") {
			url += "&" + queryStr.Encode()
		} else {
			url += "?" + queryStr.Encode()
		}
	}
	req := netlib.Get(url).SetTimeout(netlib.DEFAULT_CONNECT_TIME, api.GetTimeOut())
	return req
}

//post 请求 args[0]:map[string]string, args[1]:time.Duration
func PostForm(uri string, args ...interface{}) *netlib.HTTPRequest {
	api := GetApi(uri)
	url := api.Uri()
	req := netlib.Post(url).SetTimeout(netlib.DEFAULT_CONNECT_TIME, api.GetTimeOut())
	data := api.BandData(args...)
	if len(*data) > 0 {
		for k, v := range *data {
			req.Param(k, fmt.Sprintf("%v", v))
		}
	}
	return req
}

//post 请求 args[0]:map[string]string, args[1]:time.Duration
func PostRaw(uri string, args ...interface{}) *netlib.HTTPRequest {
	api := GetApi(uri)
	url := api.Uri()

	req := netlib.Post(url).SetTimeout(netlib.DEFAULT_CONNECT_TIME, api.GetTimeOut())
	req.JSONBody(api.BandData(args...))
	return req
}

//post 请求 args[0]:map[string]string, args[1]:time.Duration
func Put(uri string, args ...interface{}) *netlib.HTTPRequest {
	api := GetApi(uri)
	url := api.Uri()

	req := netlib.Put(url).SetTimeout(netlib.DEFAULT_CONNECT_TIME, api.GetTimeOut())
	req.JSONBody(api.BandData(args...))
	return req
}

//patch 请求 args[0]:map[string]string, args[1]:time.Duration
func Patch(uri string, args ...interface{}) *netlib.HTTPRequest {
	api := GetApi(uri)
	url := api.Uri()

	req := netlib.Patch(url).SetTimeout(netlib.DEFAULT_CONNECT_TIME, api.GetTimeOut())
	req.JSONBody(api.BandData(args...))
	return req
}

//post 请求 args[0]:map[string]string, args[1]:time.Duration
func Delete(uri string, args ...interface{}) *netlib.HTTPRequest {
	api := GetApi(uri)
	url := api.Uri()

	req := netlib.Delete(url).SetTimeout(netlib.DEFAULT_CONNECT_TIME, api.GetTimeOut())
	req.JSONBody(api.BandData(args...))
	return req
}

//delete form 请求 args[0]:map[string]string, args[1]:time.Duration
func DeleteForm(uri string, args ...interface{}) *netlib.HTTPRequest {
	api := GetApi(uri)
	url := api.Uri()
	req := netlib.Delete(url).SetTimeout(netlib.DEFAULT_CONNECT_TIME, api.GetTimeOut())

	data := api.BandData(args...)

	if len(*data) > 0 {
		for k, v := range *data {
			req.Param(k, fmt.Sprintf("%v", v))
		}
	}
	return req
}

//post 请求 args[0]:map[string]string, args[1]:time.Duration
func Head(uri string, args ...interface{}) *netlib.HTTPRequest {
	api := GetApi(uri)
	url := api.Uri()

	req := netlib.Head(url).SetTimeout(netlib.DEFAULT_CONNECT_TIME, api.GetTimeOut())
	req.JSONBody(api.BandData(args...))
	return req
}

func GetApi(uri string) (api Api) {
	uriArr := strings.Split(uri, ":")
	switch uriArr[0] {
	default:
		api = NewApiNormal(uriArr[0])
	}
	api.SetUri(uriArr[1])
	return api
}
