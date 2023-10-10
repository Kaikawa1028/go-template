//go:build wireinject

//go:generate wire gen $GOFILE

package wire

import (
	"github.com/google/wire"
	"github.com/Kaikawa1028/go-template/app/config"
	errorHandler "github.com/Kaikawa1028/go-template/app/errors/handler"

	"github.com/Kaikawa1028/go-template/app/infrastructure/persistence"
	system2 "github.com/Kaikawa1028/go-template/app/infrastructure/system"
	"github.com/Kaikawa1028/go-template/app/interfaces/handler"
	"github.com/Kaikawa1028/go-template/app/middleware"
	"github.com/Kaikawa1028/go-template/app/middleware/routeAuthorization"
	"github.com/Kaikawa1028/go-template/app/router"
	"github.com/Kaikawa1028/go-template/app/server"
	"github.com/Kaikawa1028/go-template/app/usecase"
	"github.com/Kaikawa1028/go-template/app/validate"
)

type DIContainer struct {
	Server *server.Server
	Router *router.Router
}

func InitializeDIContainer() (*DIContainer, func(), error) {
	wire.Build(
		wire.Struct(new(DIContainer), "*"),

		// Web Server(Echo)
		server.NewServer,
		router.NewRouter,
		errorHandler.NewErrorHandler,
		validate.NewCustomValidator,

		// Config
		config.NewDBConfig,
		config.NewServerConfig,

		// Middleware
		middleware.NewAccessLog,
		middleware.NewAuthenticate,
		middleware.NewAuthorize,
		middleware.NewBodyDumpLog,
		middleware.NewRecover,
		middleware.NewSetupResponseHeader,
		routeAuthorization.NewStore,
		routeAuthorization.NewCompany,

		// Handler
		handler.NewGetCompanyName,
		handler.NewGetUser,
		handler.NewHealthCheck,
		

		// Usecase
		usecase.NewAuthenticate,
		usecase.NewGetCompanyName,
		
		// Persistence
		persistence.NewOrmConnection,
		persistence.NewCompany,
		persistence.NewStoreStore,
		persistence.NewUser,
		persistence.NewAdminUserRoleOption,

		// System
		system2.NewTimer,

	)
	return &DIContainer{}, nil, nil
}
