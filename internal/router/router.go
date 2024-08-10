package router

import (
	"fmt"

	"github.com/hazkall/shadowlog/internal/controllers"
	"github.com/hazkall/shadowlog/internal/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Router struct {
	Port   int
	Server *echo.Echo
}

func NewRouter(port int) *Router {
	e := echo.New()

	e.Use(middleware.Recover())

	return &Router{
		Port:   port,
		Server: e,
	}
}

func (r *Router) Routes() *echo.Echo {
	r.Server.GET("/healthcheck", controllers.HealthCheck)

	v1 := r.Server.Group("/api/v1")
	v1.Use(middleware.Logger())
	v1.Use(middlewares.Tracing)

	{
		d := v1.Group("/deploy")
		{
			d.GET("/shoulddeploy", controllers.ShouldDeploy)
		}
	}

	return r.Server
}

func (r *Router) Start() {
	r.Server.Logger.Info("ShadowLog server is running on port", r.Port)
	r.Routes()
	if err := r.Server.Start(fmt.Sprintf("0.0.0.0:%v", r.Port)); err != nil {
		r.Server.Logger.Fatal(err)
	}
}
