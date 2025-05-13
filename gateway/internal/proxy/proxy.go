package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// func NewReverseProxy(target *url.URL) *httputil.ReverseProxy {
// 	proxy := httputil.NewSingleHostReverseProxy(target)
// 	proxy.Director = func(req *http.Request) {
// 		req.URL.Scheme = target.Scheme
// 		req.URL.Host = target.Host

// 		req.Host = target.Host
// 	}
// 	return proxy
// }

// func NewReverseProxy(target *url.URL) *httputil.ReverseProxy {
// 	proxy := httputil.NewSingleHostReverseProxy(target)
// 	proxy.Director = func(req *http.Request) {
// 		req.URL.Scheme = target.Scheme
// 		req.URL.Host = target.Host
// 		req.Host = target.Host
// 		//req.URL.Path = joinPath(target.Path, req.URL.Path) // Если нужно, добавляем префикс к пути

// 		// Log request
// 		log.Printf("Proxying request to: %s", req.URL.String())
// 	}
// 	return proxy
// }

func NewReverseProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		req.Host = target.Host // Needed for some setups
		src := req.Header

		copyHeader(req.Header, src)

		log.Printf("Proxying request to: %s", req.URL.String())
	}
	return proxy
}

func singleJoiningSlash(a, b string) string {
	aSlash := strings.HasSuffix(a, "/")
	bSlash := strings.HasPrefix(b, "/")
	switch {
	case aSlash && bSlash:
		return a + b[1:]
	case !aSlash && !bSlash:
		return a + "/" + b
	default:
		return a + b
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
