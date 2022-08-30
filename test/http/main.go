package main

import (
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// NewProxy takes target host and creates a reverse proxy
// NewProxy 拿到 targetHost 后，创建一个反向代理
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	// 请求拦截
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		modifyRequest(req, url)
	}

	// 相应拦截
	proxy.ModifyResponse = modifyResponse()
	return proxy, nil
}

// ProxyRequestHandler handles the http request using proxy
// ProxyRequestHandler 使用 proxy 处理请求
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		proxy.ServeHTTP(w, req)
	}
}

func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		resp.Header.Set("X-Proxy", "Magical")
		all, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		//fmt.Println(string(all))
		resp.Body.Close()
		resp.Body = io.NopCloser(bytes.NewBuffer(all))
		return nil
	}
}

func modifyRequest(req *http.Request, target *url.URL) {
	req.Header.Set("X-Proxy", "Simple-Reverse-Proxy")
	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host
	req.Host = target.Host //这才是关键
	req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
	if target.RawQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = target.RawQuery + req.URL.RawQuery
	} else {
		req.URL.RawQuery = target.RawQuery + "&" + req.URL.RawQuery
	}
}

func main() {
	// initialize a reverse proxy and pass the actual backend server url here
	// 初始化反向代理并传入真正后端服务的地址
	proxy, err := NewProxy("https://www.baidu.com")
	if err != nil {
		panic(err)
	}

	tlsConfig := &tls.Config{
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			return nil, nil
		},
	}

	// handle all requests to your server using the proxy
	// 使用 proxy 处理所有请求到你的服务
	http.HandleFunc("/", ProxyRequestHandler(proxy))
	//log.Fatal(http.ListenAndServeTLS(":8080", config, nil))

	httpsSrv := http.Server{
		Addr:      ":4430",
		TLSConfig: tlsConfig,
	}

	go func() {
		err := httpsSrv.ListenAndServe()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	httpSrv := http.Server{
		Addr: ":8080",
	}

	err = httpSrv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
