package server

import (
	"fmt"
	"github.com/Kaikawa1028/go-template/app/config"
	"github.com/Kaikawa1028/go-template/app/errors/handler"
	"os"
	"strings"

	"github.com/Kaikawa1028/go-template/app/middleware"
	"github.com/Kaikawa1028/go-template/app/validate"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Echo         *echo.Echo
	serverConfig *config.ServerConfig
}

func NewServer(
	errorHandler *handler.ErrorHandler,
	accessLog *middleware.AccessLog,
	recover *middleware.Recover,
	bodyDumpLog *middleware.BodyDumpLog,
	setupResponseHeader *middleware.SetupResponseHeader,
	authenticate *middleware.Authenticate,
	validator *validate.CustomValidator,
	serverConfig *config.ServerConfig,
) *Server {
	e := echo.New()

	e.HTTPErrorHandler = errorHandler.Handler
	e.Validator = validator

	e.Use(accessLog.AccessLog)
	e.Use(recover.Recover())
	e.Use(bodyDumpLog.BodyDumpLog())
	e.Use(setupResponseHeader.SetupResponseHeader)
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: strings.Split(os.Getenv("ALLOW_ORIGIN"), ","),
	}))
	e.Use(authenticate.Authenticate)

	return &Server{
		e,
		serverConfig,
	}
}

func (s Server) Start() error {
	return s.Echo.Start(fmt.Sprintf(":%d", s.serverConfig.Port))
}
