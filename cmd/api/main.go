package main

import (
	"context"
	"intellisys-backend/internal/api/httpserver"
	"intellisys-backend/internal/api/middleware"
	"intellisys-backend/internal/appfx"

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

		fx.Invoke(func(lc fx.Lifecycle, server *httpserver.Server) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go server.Start()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return server.Shutdown(ctx)
				},
			})
		}),
	).Run()
}
