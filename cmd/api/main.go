package main

import (
	"context"
	"intellisys-backend/internal/api/httpserver"
	"intellisys-backend/internal/api/middleware"
	"intellisys-backend/internal/appfx"
	"intellisys-backend/internal/domain/log"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		appfx.ConfigModule,
		appfx.PostgresModule,
		appfx.SecurityModule,
		appfx.DomainModule,
		appfx.APIModule,
		appfx.FileModlule,
		appfx.LogglerModule,

		fx.Invoke(middleware.StartRateLimitCleanup),

		fx.Invoke(func(lc fx.Lifecycle, server *httpserver.Server, log log.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go server.Start()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					log.Info("enserrando aplicação...")
					defer log.Sync()
					return server.Shutdown(ctx)
				},
			})
		}),
	).Run()
}
