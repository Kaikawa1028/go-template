//go:build wireinject && integration

//go:generate wire gen -output_file_prefix testing_ $GOFILE
package wire

import (
	"github.com/google/wire"
	"github.com/Kaikawa1028/go-template/app/config"
	"github.com/Kaikawa1028/go-template/app/domain/repository"
	"github.com/Kaikawa1028/go-template/app/domain/system"
	errorHandler "github.com/Kaikawa1028/go-template/app/errors/handler"
	"github.com/Kaikawa1028/go-template/app/infrastructure/persistence"
	"github.com/Kaikawa1028/go-template/app/middleware"
	"github.com/Kaikawa1028/go-template/app/server"
	"github.com/Kaikawa1028/go-template/app/usecase"
	"github.com/Kaikawa1028/go-template/app/validate"
	"gorm.io/gorm"
)

type IntegrationTestDIContainer struct {
	Server                        *server.Server
	OrmDB                         *gorm.DB
	Company                       repository.Company
	Store                         repository.Store
	User                          repository.User
	AdminUserRoleOption           repository.AdminUserRoleOption
	AdminUserCompany              repository.AdminUserCompany
}

func InitializeIntegrationTestDIContainer(mockTimer system.ITimer) (*IntegrationTestDIContainer, func(), error) {
	wire.Build(
		wire.Struct(new(IntegrationTestDIContainer), "*"),

		// Web Server(Echo)
		server.NewServer,
		errorHandler.NewErrorHandler,
		validate.NewCustomValidator,

		// Config
		config.NewServerConfig,

		// Middleware
		middleware.NewAccessLog,
		middleware.NewAuthenticate,
		middleware.NewBodyDumpLog,
		middleware.NewRecover,
		middleware.NewSetupResponseHeader,

		// UseCase
		usecase.NewAuthenticate,

		// Persistence
		persistence.NewOrmConnectionWithTxdb, // テスト用のDB接続(gorm)
		persistence.NewCompany,
		persistence.NewStoreStore,
		persistence.NewUser,
		persistence.NewAdminUserRoleOption,
		persistence.NewAdminUserCompany,


		// Externals
		// ExternalsはDIしない
		// テスト時に手動でモックを注入する

		// Resources
		//testResources.NewResources,
	)
	return &IntegrationTestDIContainer{}, nil, nil
}
