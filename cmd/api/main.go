package main

import (
	"financas/configuration"
	"financas/internal/api"
	"financas/internal/config"
)

func main() {
	c := configuration.New()
	var cfg config.Config

	cfg.Port = c.Server.Port
	cfg.Env = "development"
	cfg.DB.DSN = c.DB.DSN
	cfg.DB.MaxOpenConns = c.DB.MaxOpenConns
	cfg.DB.MaxIdleConns = c.DB.MaxIdleConns
	cfg.DB.MaxIdleTime = c.DB.MaxIdleTime
	cfg.Limiter.RPS = c.RateLimiter.RPS
	cfg.Limiter.Burst = c.RateLimiter.Burst
	cfg.Limiter.Enabled = c.RateLimiter.Enabled

	app := api.NewApp(cfg)
	err := app.Serve()
	if err != nil {
		app.Logger.PrintFatal(err, nil)
	}

}
