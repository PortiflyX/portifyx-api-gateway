package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/portifyx/portifyx-api-gateway/config"
	"github.com/portifyx/portifyx-api-gateway/internal/rest"
	"github.com/portifyx/portifyx-goutils/pkg/logx"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start HTTP edge reverse proxy server",
	Run:   runAPIServer,
}

func init() {
	rootCmd.AddCommand(apiCmd)
}

func runAPIServer(cmd *cobra.Command, args []string) {
	config.Init()
	cfg := config.GetConfig()

	logger, err := logx.Init(cfg.App.AppName, cfg.App.AppVersion, cfg.App.AppEnv)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	ctx := context.Background()
	logx.Infof(ctx, "[*] starting %s v%s on port :%s (env: %s)", cfg.App.AppName, cfg.App.AppVersion, cfg.App.AppPort, cfg.App.AppEnv)

	srv := rest.NewEchoServer(cfg)
	e := srv.GetEcho()

	go func() {
		addr := fmt.Sprintf(":%s", cfg.App.AppPort)
		if err := e.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logx.Fatalf(ctx, err, "[!] fatal error starting echo server: %v", err)
		}
	}()

	gracefulShutdown(e)
}
