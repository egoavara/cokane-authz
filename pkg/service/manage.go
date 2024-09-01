package service

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"egoavara.net/authz/pkg/util"
	"github.com/gin-gonic/gin"
)

type Manage struct {
	proxyUrl   *url.URL
	proxy      *httputil.ReverseProxy
	Prometheus *PrometheusExporter
}

func (s *Manage) Setup(engine *gin.Engine) {
	api := engine.Group("/api/:version/")
	api.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": fmt.Sprintf("Hello from %s", context.Param("version")),
		})
	})
	api.Any("/openfga/*paths", func(context *gin.Context) {
		s.proxy.Director = func(req *http.Request) {
			req.Header = context.Request.Header
			req.Host = s.proxyUrl.Host
			req.URL.Scheme = s.proxyUrl.Scheme
			req.URL.Host = s.proxyUrl.Host
			req.URL.Path = context.Param("paths")
			log.Println("Reverse Proxy OpenFGA to", req.URL)
		}
		s.proxy.ServeHTTP(context.Writer, context.Request)
	})
	s.Prometheus.Setup(engine, []string{"/meta/:version/metrics", "/metrics"})
}

func NewManage() *Manage {
	proxyUrl := util.Must(url.Parse("http://localhost:8080"))
	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	Prometheus := NewPrometheusExporter("cokane_authz", "manage")
	return &Manage{
		proxyUrl,
		proxy,
		Prometheus,
	}
}
