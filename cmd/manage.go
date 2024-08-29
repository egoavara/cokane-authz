/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
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
				"message": "Not Hello World",
			})
		})
		engine.Run(":80")
	},
}

func init() {
	rootCmd.AddCommand(manageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// manageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// manageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
