package services

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ProxyService struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func NewProxyService(targetURL string) *ProxyService {
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("Failed to parse target URL: %v", err)
	}

	return &ProxyService{
		target: target,
		proxy:  httputil.NewSingleHostReverseProxy(target),
	}
}

func (p *ProxyService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}