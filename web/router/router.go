package router

import (
	"github.com/labstack/echo/v4"
	"github.com/Kaikawa1028/go-template/app/const/app"
	"github.com/Kaikawa1028/go-template/app/interfaces/handler"
	"github.com/Kaikawa1028/go-template/app/middleware"
	"github.com/Kaikawa1028/go-template/app/middleware/routeAuthorization"
)

type Router struct {
	authorize                        *middleware.Authorize
	routeAuthStore                   *routeAuthorization.Store
	routeAuthCompany                 *routeAuthorization.Company

	getCompanyName                                *handler.GetCompanyName
	getUser                                       *handler.GetUser
	healthCheck                                   *handler.HealthCheck
	
}

func NewRouter(
	authorize *middleware.Authorize,
	routeAuthStore *routeAuthorization.Store,
	routeAuthCompany *routeAuthorization.Company,

	getCompanyName *handler.GetCompanyName,
	getUser *handler.GetUser,
	healthCheck *handler.HealthCheck,
	
) *Router {
	return &Router{
		authorize,
		routeAuthStore,
		routeAuthCompany,
		getCompanyName,
		getUser,
		healthCheck,
	}
}

func (r Router) Attach(e *echo.Echo) {
	g := e.Group("")



	group(g, "/v1/companies/:company-id", []echo.MiddlewareFunc{r.routeAuthCompany.AuthorizeCompany}, func(g *echo.Group) {

		g.GET("/name", r.getCompanyName.GetCompanyName, r.authorize.Authorize([]string{app.RoleAgency, app.RoleManager, app.RoleEditor, app.RoleWorker, app.RoleViewer}))
	})

	g.GET("/v1/healthcheck", r.healthCheck.HealthCheck)

	g.GET("/v1/user", r.getUser.GetUser)
}

func group(g *echo.Group, prefix string, m []echo.MiddlewareFunc, cb func(*echo.Group)) {
	childGroup := g.Group(prefix, m...)
	cb(childGroup)
}
