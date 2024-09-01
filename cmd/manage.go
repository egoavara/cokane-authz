/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"egoavara.net/authz/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	PROXY_URL *url.URL
	PROXY     *httputil.ReverseProxy
)

// manageCmd represents the manage command
var manageCmd = &cobra.Command{
	Use:   "manage",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		engine := gin.Default()
		engine.GET("/", func(context *gin.Context) {
			context.JSON(200, gin.H{
				"message": "Hello World 2024-09-01T16:37:00",
			})
		})
		engine.Any("/stores/*paths", func(context *gin.Context) {

			PROXY.Director = func(req *http.Request) {
				req.Header = context.Request.Header
				req.Host = PROXY_URL.Host
				req.URL.Scheme = PROXY_URL.Scheme
				req.URL.Host = PROXY_URL.Host
				req.URL.Path = context.Param("paths")
				log.Println("req", req.URL)
			}
			log.Println("Proxying to", PROXY_URL)
			PROXY.ServeHTTP(context.Writer, context.Request)
		})
		engine.Run(":80")
	},
}

func init() {
	rootCmd.AddCommand(manageCmd)

	PROXY_URL = util.Must(url.Parse("http://localhost:8080"))
	PROXY = httputil.NewSingleHostReverseProxy(PROXY_URL)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// manageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// manageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
