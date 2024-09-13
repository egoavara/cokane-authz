/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"net/http"
	"time"

	"egoavara.net/authz/pkg/manage"
	"github.com/go-faster/sdk/app"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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
		var shutdownTimeout = 15 * time.Second
		app.Run(func(ctx context.Context, log *zap.Logger, m *app.Metrics) error {
			router, err := manage.NewServer(
				manage.Entrypoint{},
				manage.WithMeterProvider(m.MeterProvider()),
				manage.WithTracerProvider(m.TracerProvider()),
			)
			if err != nil {
				return err
			}
			server := http.Server{
				ReadTimeout: 5 * time.Second,
				Addr:        ":80",
				Handler:     router,
			}
			g, ctx := errgroup.WithContext(ctx)
			g.Go(func() error {
				// Wait until g ctx canceled, then try to shut down server.
				<-ctx.Done()

				log.Info("Shutting down", zap.Duration("timeout", shutdownTimeout))

				shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
				defer cancel()
				return server.Shutdown(shutdownCtx)
			})
			g.Go(func() error {
				defer log.Info("Server stopped")
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					return err
				}
				return nil
			})

			return g.Wait()
		})
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
