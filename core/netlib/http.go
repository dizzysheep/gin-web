package netlib

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gin-web/core/config"
	"gin-web/core/log"
	"gin-web/core/skywalking"
	"github.com/gin-gonic/gin"
	"github.com/rs/dnscache"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/objx"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var defaultSetting = HttpSettings{
	UserAgent:        "server",
	ConnectTimeout:   60 * time.Second,
	ReadWriteTimeout: 60 * time.Second,
	Gzip:             true,
	DumpBody:         true,
}

// http默认连接超时时间
const DefaultConnectTime = time.Second * 3

var defaultCookieJar http.CookieJar
var settingMutex sync.Mutex

// 定义全局的 DNS cache
var cacheDNS *dnscache.Resolver

// 定义 DNS cache 过期时间
const cacheDNSExpire = 5 * time.Minute

// 是否默认记录请求日志(默认开启)
var requestLogOn = true

func init() {
	// 初始化 DNS cache
	cacheDNS = &dnscache.Resolver{}

	//每5分钟更新一下dns缓存
	go func() {
		t := time.NewTicker(cacheDNSExpire)
		defer t.Stop()
		for range t.C {
			cacheDNS.Refresh(true)
		}
	}()

	//是否设置了关闭默认请求日志
	apisRequestLogOff := config.GetInt("log.apisRequestLogOff")
	if apisRequestLogOff == 1 {
		requestLogOn = false
	}

}

// createDefaultCookie creates a global cookiejar to store cookies.
func createDefaultCookie() {
	settingMutex.Lock()
	defer settingMutex.Unlock()
	defaultCookieJar, _ = cookiejar.New(nil)
}

// SetDefaultSetting Overwrite default settings
func SetDefaultSetting(setting HttpSettings) {
	settingMutex.Lock()
	defer settingMutex.Unlock()
	defaultSetting = setting
}

// 生成一串唯一识别码字符串
func generatuniqueCode() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < 20; i++ {
		result = append(result, Bytes[rand.Intn(len(Bytes))])
	}
	return string(result)
}

// NewRequest return *HTTPRequest with specific method
func NewRequest(rawurl, method string) *HTTPRequest {
	var resp http.Response
	u, err := url.Parse(rawurl)
	if err != nil {
		log.Loger.Println("Httplib:", err)
	}
	req := http.Request{
		URL:        u,
		Method:     method,
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	return &HTTPRequest{
		url:        rawurl,
		req:        &req,
		params:     map[string][]string{},
		files:      map[string]string{},
		setting:    defaultSetting,
		resp:       &resp,
		uniqueCode: generatuniqueCode(),
	}
}

// Get returns *HTTPRequest with GET method.
func Get(url string) *HTTPRequest {
	return NewRequest(url, "GET")
}

// Post returns *HTTPRequest with POST method.
func Post(url string) *HTTPRequest {
	return NewRequest(url, "POST")
}

// Put returns *HTTPRequest with PUT method.
func Put(url string) *HTTPRequest {
	return NewRequest(url, "PUT")
}

// Patch returns *HTTPRequest with PUT method.
func Patch(url string) *HTTPRequest {
	return NewRequest(url, "PATCH")
}

// Delete returns *HTTPRequest DELETE method.
func Delete(url string) *HTTPRequest {
	return NewRequest(url, "DELETE")
}

// Head returns *HTTPRequest with HEAD method.
func Head(url string) *HTTPRequest {
	return NewRequest(url, "HEAD")
}

// HttpSettings is the http.Client setting
type HttpSettings struct {
	ShowDebug        bool
	UserAgent        string
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
	TLSClientConfig  *tls.Config
	Proxy            func(*http.Request) (*url.URL, error)
	Transport        http.RoundTripper
	CheckRedirect    func(req *http.Request, via []*http.Request) error
	EnableCookie     bool
	Gzip             bool
	DumpBody         bool
	Retries          int // if set to -1 means will retry forever
}

// HTTPRequest provides more useful methods for requesting one url than http.Request.
type HTTPRequest struct {
	url        string
	params     map[string][]string
	req        *http.Request
	files      map[string]string
	setting    HttpSettings
	resp       *http.Response
	body       []byte
	dump       []byte
	ctx        *gin.Context //用户框架上下文，用于从中读取skywalking上下文
	uniqueCode string       //当前请求唯一识别码
}

func (b *HTTPRequest) SetCtx(ctx *gin.Context) *HTTPRequest {
	b.ctx = ctx
	return b
}

// get Url with no get params
func (b *HTTPRequest) GetRequestUrlWithNoGetParams() string {
	reqUrl := b.url
	reqs := strings.Split(reqUrl, "?")
	return reqs[0]
}

// GetRequest return the request object
func (b *HTTPRequest) GetRequest() *http.Request {
	return b.req
}

// Setting Change request settings
func (b *HTTPRequest) Setting(setting HttpSettings) *HTTPRequest {
	b.setting = setting
	return b
}

// SetBasicAuth sets the request's Authorization header to use HTTP Basic Authentication with the provided username and password.
func (b *HTTPRequest) SetBasicAuth(username, password string) *HTTPRequest {
	b.req.SetBasicAuth(username, password)
	return b
}

// SetEnableCookie sets enable/disable cookiejar
func (b *HTTPRequest) SetEnableCookie(enable bool) *HTTPRequest {
	b.setting.EnableCookie = enable
	return b
}

// SetUserAgent sets User-Agent header field
func (b *HTTPRequest) SetUserAgent(useragent string) *HTTPRequest {
	b.setting.UserAgent = useragent
	return b
}

// Debug sets show debug or not when executing request.
func (b *HTTPRequest) Debug(isdebug bool) *HTTPRequest {
	b.setting.ShowDebug = isdebug
	return b
}

// Retries sets Retries times.
// default is 0 means no retried.
// -1 means retried forever.
// others means retried times.
func (b *HTTPRequest) Retries(times int) *HTTPRequest {
	b.setting.Retries = times
	return b
}

// DumpBody setting whether need to Dump the Body.
func (b *HTTPRequest) DumpBody(isdump bool) *HTTPRequest {
	b.setting.DumpBody = isdump
	return b
}

// DumpRequest return the DumpRequest
func (b *HTTPRequest) DumpRequest() []byte {
	return b.dump
}

// SetTimeout sets connect time out and read-write time out for HTTPRequest.
func (b *HTTPRequest) SetTimeout(connectTimeout, readWriteTimeout time.Duration) *HTTPRequest {
	b.setting.ConnectTimeout = connectTimeout
	b.setting.ReadWriteTimeout = readWriteTimeout
	return b
}

// SetTLSClientConfig sets tls connection configurations if visiting https url.
func (b *HTTPRequest) SetTLSClientConfig(config *tls.Config) *HTTPRequest {
	b.setting.TLSClientConfig = config
	return b
}

// SetContext set the request context.
func (b *HTTPRequest) SetContext(ctx context.Context) *HTTPRequest {
	b.req = b.req.WithContext(ctx)
	return b
}

// Header add header item string in request.
func (b *HTTPRequest) Header(key, value string) *HTTPRequest {
	b.req.Header.Set(key, value)
	return b
}

// SetHost set the request host
func (b *HTTPRequest) SetHost(host string) *HTTPRequest {
	b.req.Host = host
	return b
}

// SetProtocolVersion Set the protocol version for incoming requests.
// Client requests always use HTTP/1.1.
func (b *HTTPRequest) SetProtocolVersion(vers string) *HTTPRequest {
	if len(vers) == 0 {
		vers = "HTTP/1.1"
	}

	major, minor, ok := http.ParseHTTPVersion(vers)
	if ok {
		b.req.Proto = vers
		b.req.ProtoMajor = major
		b.req.ProtoMinor = minor
	}

	return b
}

// SetCookie add cookie into request.
func (b *HTTPRequest) SetCookie(cookie *http.Cookie) *HTTPRequest {
	b.req.Header.Add("Cookie", cookie.String())
	return b
}

// SetTransport set the setting transport
func (b *HTTPRequest) SetTransport(transport http.RoundTripper) *HTTPRequest {
	b.setting.Transport = transport
	return b
}

// SetProxy set the http proxy
// example:
//
//	func(req *http.Request) (*url.URL, error) {
//		u, _ := url.ParseRequestURI("http://127.0.0.1:8118")
//		return u, nil
//	}
func (b *HTTPRequest) SetProxy(proxy func(*http.Request) (*url.URL, error)) *HTTPRequest {
	b.setting.Proxy = proxy
	return b
}

// SetCheckRedirect specifies the policy for handling redirects.
//
// If CheckRedirect is nil, the Client uses its default policy,
// which is to stop after 10 consecutive requests.
func (b *HTTPRequest) SetCheckRedirect(redirect func(req *http.Request, via []*http.Request) error) *HTTPRequest {
	b.setting.CheckRedirect = redirect
	return b
}

// Param adds query param in to request.
// params build query string as ?key1=value1&key2=value2...
func (b *HTTPRequest) Param(key, value string) *HTTPRequest {
	if param, ok := b.params[key]; ok {
		b.params[key] = append(param, value)
	} else {
		b.params[key] = []string{value}
	}
	return b
}

// PostFile add a post file to the request
func (b *HTTPRequest) PostFile(formname, filename string) *HTTPRequest {
	b.files[formname] = filename
	return b
}

// Body adds request raw body.
// it supports string and []byte.
func (b *HTTPRequest) Body(data interface{}) *HTTPRequest {
	switch t := data.(type) {
	case string:
		bf := bytes.NewBufferString(t)
		b.req.Body = ioutil.NopCloser(bf)
		b.req.ContentLength = int64(len(t))
	case []byte:
		bf := bytes.NewBuffer(t)
		b.req.Body = ioutil.NopCloser(bf)
		b.req.ContentLength = int64(len(t))
	}
	return b
}

// XMLBody adds request raw body encoding by XML.
func (b *HTTPRequest) XMLBody(obj interface{}) (*HTTPRequest, error) {
	if b.req.Body == nil && obj != nil {
		byts, err := xml.Marshal(obj)
		if err != nil {
			return b, err
		}
		b.req.Body = ioutil.NopCloser(bytes.NewReader(byts))
		b.req.ContentLength = int64(len(byts))
		b.req.Header.Set("Content-Type", "application/xml")
	}
	return b, nil
}

// YAMLBody adds request raw body encoding by YAML.
func (b *HTTPRequest) YAMLBody(obj interface{}) (*HTTPRequest, error) {
	if b.req.Body == nil && obj != nil {
		byts, err := yaml.Marshal(obj)
		if err != nil {
			return b, err
		}
		b.req.Body = ioutil.NopCloser(bytes.NewReader(byts))
		b.req.ContentLength = int64(len(byts))
		b.req.Header.Set("Content-Type", "application/x+yaml")
	}
	return b, nil
}

// JSONBody adds request raw body encoding by JSON.
func (b *HTTPRequest) JSONBody(obj interface{}) (*HTTPRequest, error) {
	if b.req.Body == nil && obj != nil {
		byts, err := json.Marshal(obj)
		if err != nil {
			return b, err
		}
		b.req.Body = ioutil.NopCloser(bytes.NewReader(byts))
		b.req.ContentLength = int64(len(byts))
		b.req.Header.Set("Content-Type", "application/json")
	}
	return b, nil
}

func (b *HTTPRequest) buildURL(paramBody string) {
	// build GET url with query string
	if b.req.Method == "GET" && len(paramBody) > 0 {
		if strings.Contains(b.url, "?") {
			b.url += "&" + paramBody
		} else {
			b.url = b.url + "?" + paramBody
		}
		return
	}

	// build POST/PUT/PATCH url and body
	if (b.req.Method == "POST" || b.req.Method == "PUT" || b.req.Method == "PATCH" || b.req.Method == "DELETE") && b.req.Body == nil {
		// with files
		if len(b.files) > 0 {
			pr, pw := io.Pipe()
			bodyWriter := multipart.NewWriter(pw)
			go func() {
				for formname, filename := range b.files {
					fileWriter, err := bodyWriter.CreateFormFile(formname, filename)
					if err != nil {
						log.Loger.Println("Httplib:", err)
					}
					fh, err := os.Open(filename)
					if err != nil {
						log.Loger.Println("Httplib:", err)
					}
					//iocopy
					_, err = io.Copy(fileWriter, fh)
					fh.Close()
					if err != nil {
						log.Loger.Println("Httplib:", err)
					}
				}
				for k, v := range b.params {
					for _, vv := range v {
						bodyWriter.WriteField(k, vv)
					}
				}
				bodyWriter.Close()
				pw.Close()
			}()
			b.Header("Content-Type", bodyWriter.FormDataContentType())
			b.req.Body = ioutil.NopCloser(pr)
			return
		}

		// with params
		if len(paramBody) > 0 {
			b.Header("Content-Type", "application/x-www-form-urlencoded")
			b.Body(paramBody)
		}
	}
}

func (b *HTTPRequest) getResponse() (*http.Response, error) {
	if b.resp.StatusCode != 0 {
		return b.resp, nil
	}

	//接入skywalking，向Header注入sw追踪信息
	if config.SkywalkingSwitch {
		reqUrl := b.GetRequestUrlWithNoGetParams()
		exitSpan, headers := skywalking.CreateHttpSwExitSpan(reqUrl, b.req.Method, b.ctx.Request.Context())
		for _, header := range headers {
			b.Header(header[0], header[1])
		}
		defer func() {
			if exitSpan != nil {
				exitSpan.End()
			}
		}()
	}

	//header增加唯一识别码
	b.Header("uniqueCode", b.uniqueCode)

	resp, err := b.DoRequest()
	if err != nil {
		return nil, err
	}
	b.resp = resp
	return resp, nil
}

// DoRequest will do the client.Do
func (b *HTTPRequest) DoRequest() (resp *http.Response, err error) {
	var paramBody string
	if len(b.params) > 0 {
		var buf bytes.Buffer
		for k, v := range b.params {
			for _, vv := range v {
				buf.WriteString(url.QueryEscape(k))
				buf.WriteByte('=')
				buf.WriteString(url.QueryEscape(vv))
				buf.WriteByte('&')
			}
		}
		paramBody = buf.String()
		paramBody = paramBody[0 : len(paramBody)-1]
	}

	b.buildURL(paramBody)
	urlParsed, err := url.Parse(b.url)
	if err != nil {
		return nil, err
	}

	b.req.URL = urlParsed

	trans := b.setting.Transport

	if trans == nil {
		// create default transport
		trans = &http.Transport{
			TLSClientConfig:     b.setting.TLSClientConfig,
			Proxy:               b.setting.Proxy,
			MaxIdleConnsPerHost: 100,
			DialContext:         TimeoutDialer(b.setting.ConnectTimeout, b.setting.ReadWriteTimeout),
		}
	} else {
		// if b.transport is *http.Transport then set the settings.
		if t, ok := trans.(*http.Transport); ok {
			if t.TLSClientConfig == nil {
				t.TLSClientConfig = b.setting.TLSClientConfig
			}
			if t.Proxy == nil {
				t.Proxy = b.setting.Proxy
			}
			if t.DialContext == nil {
				t.DialContext = TimeoutDialer(b.setting.ConnectTimeout, b.setting.ReadWriteTimeout)
			}
		}
	}

	var jar http.CookieJar
	if b.setting.EnableCookie {
		if defaultCookieJar == nil {
			createDefaultCookie()
		}
		jar = defaultCookieJar
	}

	client := &http.Client{
		Transport: trans,
		Jar:       jar,
	}

	if b.setting.UserAgent != "" && b.req.Header.Get("User-Agent") == "" {
		b.req.Header.Set("User-Agent", b.setting.UserAgent)
	}

	if b.setting.CheckRedirect != nil {
		client.CheckRedirect = b.setting.CheckRedirect
	}

	if b.setting.ShowDebug {
		dump, err := httputil.DumpRequest(b.req, b.setting.DumpBody)
		if err != nil {
			log.Get(b.ctx).Error(err.Error())
		}
		b.dump = dump
	}

	// retries default value is 0, it will run once.
	// retries equal to -1, it will run forever until success
	// retries is setted, it will retries fixed times.
	for i := 0; b.setting.Retries == -1 || i <= b.setting.Retries; i++ {
		resp, err = client.Do(b.req)
		if err == nil {
			break
		}
	}

	return resp, err
}

// String returns the body string in response.
// it calls Response inner.
func (b *HTTPRequest) String() (string, error) {
	data, err := b.Bytes()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Bytes returns the body []byte in response.
// it calls Response inner.
func (b *HTTPRequest) Bytes() ([]byte, error) {
	if b.body != nil {
		return b.body, nil
	}
	resp, err := b.getResponse()
	if err != nil {
		return nil, err
	}
	if resp.Body == nil {
		return nil, nil
	}
	defer resp.Body.Close()
	if b.setting.Gzip && resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		b.body, err = ioutil.ReadAll(reader)
	} else {
		b.body, err = ioutil.ReadAll(resp.Body)
	}

	//默认记录请求日志，默认开启，可通过在配置中设置apisRequestLogOff=1关闭
	if requestLogOn {
		log.Get(b.ctx).WithFields(logrus.Fields{
			"api_url":      b.url,
			"api_params":   b.params,
			"api_response": string(b.body),
		}).Info("api info")
	}

	return b.body, err
}

// ToFile saves the body data in response to one file.
// it calls Response inner.
func (b *HTTPRequest) ToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	resp, err := b.getResponse()
	if err != nil {
		return err
	}
	if resp.Body == nil {
		return nil
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// ToJSON returns the map that marshals from the body bytes as json in response .
// it calls Response inner.
func (b *HTTPRequest) ToJSON(v interface{}) error {
	startTime := time.Now().UnixNano()
	data, err := b.Bytes()
	if err != nil {
		return err
	}

	//添加响应时间
	timeGap := time.Now().UnixNano() - startTime
	timeGapFloat := float64(timeGap) / 1000000000
	timeGapString := strconv.FormatFloat(timeGapFloat, 'f', 3, 64)

	err = json.Unmarshal(data, v)
	if err != nil {
		msg := string(data)
		if len(msg) > 5000 {
			msg = msg[:5000]
		}
		errMsg := fmt.Sprintf(
			`netlib.ToJSON解析错误，请求url：%s
响应原文：%s
错误信息：%s
唯一识别码(REQUEST HEADER: uniqueCode)：%s
请求耗时：%s秒`,
			b.url,
			msg,
			err.Error(),
			b.uniqueCode,
			timeGapString,
		)
		log.Loger.Error(errMsg)
	}
	return err
}

// ToMap returns the map that marshals from the body bytes as json in response .
// it calls Response inner.
func (b *HTTPRequest) ToMap() (objx.Map, error) {
	data, err := b.Bytes()
	if err != nil {
		return nil, err
	}
	return objx.FromJSON(string(data))
}

// ToXML returns the map that marshals from the body bytes as xml in response .
// it calls Response inner.
func (b *HTTPRequest) ToXML(v interface{}) error {
	data, err := b.Bytes()
	if err != nil {
		return err
	}
	return xml.Unmarshal(data, v)
}

// ToYAML returns the map that marshals from the body bytes as yaml in response .
// it calls Response inner.
func (b *HTTPRequest) ToYAML(v interface{}) error {
	data, err := b.Bytes()
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, v)
}

// Response executes request client gets response mannually.
func (b *HTTPRequest) Response() (*http.Response, error) {
	return b.getResponse()
}

// TimeoutDialer returns functions of connection dialer with timeout settings for http.Transport Dial field.
func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(ctx context.Context, net, addr string) (c net.Conn, err error) {
	return func(ctx context.Context, netw, addr string) (net.Conn, error) {
		// 通过请求url获取域名和port
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		ips, err := cacheDNS.LookupHost(ctx, host)

		if err != nil {
			log.Loger.Errorf("cacheDNS.LookupHost %s", addr)
			return nil, err
		}

		var conn net.Conn
		for _, ip := range ips {
			dialer := net.Dialer{
				Timeout:   cTimeout,
				KeepAlive: 30 * time.Second,
			}
			conn, err = dialer.DialContext(ctx, netw, net.JoinHostPort(ip, port))
			if err == nil {
				//设置读写超时时间
				conn.SetDeadline(time.Now().Add(rwTimeout))
				break
			}
		}

		return conn, err
	}
}
