package netlib

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gin-web/core/config"
	"github.com/google/uuid"
	"github.com/rs/dnscache"
	"io"
	"io/ioutil"
	"log"
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

// 是否默认记录请求日志(默认开启)
var requestLogOn = true

var defaultCookieJar http.CookieJar
var settingMutex sync.Mutex

var cacheDNS *dnscache.Resolver

//DNS缓存时间
const cacheDNSExpire = 5 * time.Minute

// http默认连接超时时间
const DEFAULT_CONNECT_TIME = time.Second * 3

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

var defaultSetting = HttpSettings{
	UserAgent:        "server",
	ConnectTimeout:   60 * time.Second,
	ReadWriteTimeout: 60 * time.Second,
	Gzip:             true,
	DumpBody:         true,
}

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

// 上下文对象，兼容gin框架上下文gin.Context
type swContext interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	GetHeader(headerKey string) string
}

type HTTPRequest struct {
	url           string
	params        map[string][]string
	req           *http.Request
	files         map[string]string
	setting       HttpSettings
	resp          *http.Response
	body          []byte
	dump          []byte
	skywalkingOn  bool      //是否设置skywalking链路追踪
	skywalkingCtx swContext //用户框架上下文，用于从中读取skywalking上下文
	uniqueCode    string    //当前请求唯一识别码
}

//@desc NewRequest 创建request
//rawurl请求地址
func NewRequest(rawurl, method string) *HTTPRequest {
	var resp http.Response
	u, err := url.Parse(rawurl)
	if err != nil {
		log.Println("Httplib:", err)
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
		uniqueCode: getUniqueNo(),
	}
}

func (b *HTTPRequest) getRequest() *http.Request {
	return b.req
}

// Debug 是否使用debug
func (b *HTTPRequest) Debug(isdebug bool) *HTTPRequest {
	b.setting.ShowDebug = isdebug
	return b
}

// Retries 设置重试次数
func (b *HTTPRequest) Retries(times int) *HTTPRequest {
	b.setting.Retries = times
	return b
}

// SetTimeout 设置超时时间
func (b *HTTPRequest) SetTimeout(connectTimeout, readWriteTimeout time.Duration) *HTTPRequest {
	b.setting.ConnectTimeout = connectTimeout
	b.setting.ReadWriteTimeout = readWriteTimeout
	return b
}

func (b *HTTPRequest) Header(key, value string) *HTTPRequest {
	b.req.Header.Set(key, value)
	return b
}

// Param 添加参数
func (b *HTTPRequest) Param(key, value string) *HTTPRequest {
	if param, ok := b.params[key]; ok {
		b.params[key] = append(param, value)
	} else {
		b.params[key] = []string{value}
	}
	return b
}

// Body 直接信息返回到Body
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
		log.Println(errMsg)
	}
	return err
}

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
		log.Println("ApisRequestRecord ●url:%+v  ●param:%+v  ●response:%+v", b.url, b.params, string(b.body))
	}

	return b.body, err
}

func (b *HTTPRequest) getResponse() (*http.Response, error) {
	if b.resp.StatusCode != 0 {
		return b.resp, nil
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
			log.Println(err.Error())
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
						log.Println("Httplib:", err)
					}
					fh, err := os.Open(filename)
					if err != nil {
						log.Println("Httplib:", err)
					}
					//iocopy
					_, err = io.Copy(fileWriter, fh)
					fh.Close()
					if err != nil {
						log.Println("Httplib:", err)
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
			log.Printf("cacheDNS.LookupHost %s", addr)
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

//获取唯一请求code
func getUniqueNo() string {
	return uuid.New().String()
}

// Get returns *HttpRequest with GET method.
func Get(url string) *HTTPRequest {
	return NewRequest(url, "GET")
}

// Post returns *HttpRequest with POST method.
func Post(url string) *HTTPRequest {
	return NewRequest(url, "POST")
}

// Put returns *HttpRequest with PUT method.
func Put(url string) *HTTPRequest {
	return NewRequest(url, "PUT")
}

// Patch returns *HttpRequest with PUT method.
func Patch(url string) *HTTPRequest {
	return NewRequest(url, "PATCH")
}

// Delete returns *HttpRequest DELETE method.
func Delete(url string) *HTTPRequest {
	return NewRequest(url, "DELETE")
}

// Head returns *HttpRequest with HEAD method.
func Head(url string) *HTTPRequest {
	return NewRequest(url, "HEAD")
}

func createDefaultCookie() {
	settingMutex.Lock()
	defer settingMutex.Unlock()
	defaultCookieJar, _ = cookiejar.New(nil)
}
