package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func Proxy(c *gin.Context, host, path string) {
	u, _ := url.ParseRequestURI(host)
	u.Path = path
	targetQuery := u.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = u.Scheme
		req.URL.Host = u.Host
		req.URL.Path = u.Path
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}

	reverseProxy := &httputil.ReverseProxy{Director: director}
	reverseProxy.ServeHTTP(c.Writer, c.Request)
	return
}
