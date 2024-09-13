package controller

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

type OpenfgaProxy struct {
	proxyUrl *url.URL
	proxy    *httputil.ReverseProxy
}

func NewOpenfgaProxy(addr string) (*OpenfgaProxy, error) {
	proxyUrl, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	return &OpenfgaProxy{
		proxyUrl,
		proxy,
	}, nil
}

func (s *OpenfgaProxy) UseRoute(engine *gin.Engine, router gin.IRouter) {
	router.Any("/*paths", func(c *gin.Context) {
		s.proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = s.proxyUrl.Host
			req.URL.Scheme = s.proxyUrl.Scheme
			req.URL.Host = s.proxyUrl.Host
			req.URL.Path = c.Param("paths")
			log.Println("Reverse Proxy OpenFGA to", req.URL)
		}
		s.proxy.ServeHTTP(c.Writer, c.Request)
	})
}
