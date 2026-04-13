package appfx

import (
	"intellisys-backend/internal/api/controller"
	"intellisys-backend/internal/api/httpserver"
	"intellisys-backend/internal/api/middleware"
	"intellisys-backend/internal/domain/log"
	"intellisys-backend/internal/domain/shortcut"
	"intellisys-backend/internal/domain/user"
	"intellisys-backend/internal/infra/config"
	"intellisys-backend/internal/infra/logger"
	"intellisys-backend/internal/infra/persistence/postgres"
	"intellisys-backend/internal/infra/security"
	"intellisys-backend/internal/infra/storage"
	"net/http"

	"go.uber.org/fx"
)

var ConfigModule = fx.Module("config",
	fx.Provide(config.LoadConfig),
)

var PostgresModule = fx.Module("postgres",
	fx.Provide(
		postgres.NewPool,
		fx.Annotate(postgres.NewUserRepository, fx.As(new(user.UserRepository))),
		fx.Annotate(postgres.NewShortcutRepository, fx.As(new(shortcut.ShortcutRepository))),
	),
)

var SecurityModule = fx.Module("security",
	fx.Provide(
		fx.Annotate(security.NewBcryptHasher, fx.As(new(user.PasswordHasher))),
		security.NewJWTService,
		func(j *security.JWTService) user.TokenGenerator { return j },
	),
)

var DomainModule = fx.Module("domain",
	fx.Provide(
		shortcut.NewshortcutService,
		user.NewUserService,
		user.NewAuthService,
	),
)

var APIModule = fx.Module("api",
	fx.Provide(
		controller.NewshortcutController,
		controller.NewUserController,
		controller.NewAuthController,
		middleware.NewAuthMidware,
		fx.Annotate(
			httpserver.NewRouter,
			fx.As(new(http.Handler)),
			fx.ResultTags(
				``, // ShortcutController
				``, // UserController
				``, // AuthController
				``, // AuthMidware
				``, // Config
				``, // Logger
				`name:"globalLimiter"`,
				`name:"loginLimiter"`,
			),
		),
		httpserver.NewServer,
	),
)

var FileModlule = fx.Module("file",
	fx.Provide(
		fx.Annotate(storage.NewLocalStorage, fx.As(new(shortcut.FileStorage))),
	),
)

var LogglerModule = fx.Module("logger",
	fx.Provide(
		fx.Annotate(logger.NewLogger, fx.As(new(log.Logger))),
	),
)

var RateLimitModule = fx.Module("ratelimit",
	fx.Provide(
		fx.Annotate(
			middleware.NewGlobalRateLimiter,
			fx.ResultTags(`name:"globalLimiter"`),
		),
		fx.Annotate(
			middleware.NewLoginRateLimiter,
			fx.ResultTags(`name:"loginLimiter"`),
		),
	),
)
