package web

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/murtaza-u/rinc/internal/conf"
	"github.com/murtaza-u/rinc/internal/util"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type Srv struct {
	conf   conf.C
	router *echo.Echo
}

func NewSrv(c conf.C) Srv {
	r := echo.New()
	r.Pre(echoMiddleware.RemoveTrailingSlash()) // trim trailing slash
	return Srv{
		conf:   c,
		router: r,
	}
}

func (s Srv) Run() error {
	slog.SetDefault(util.NewLogger(s.conf.Log))
	err := os.MkdirAll(s.conf.Output, 0o755)
	if err != nil {
		return fmt.Errorf("creating %q directory: %w", s.conf.Output, err)
	}
	s.router.Static("/static", filepath.Join("static"))
	s.router.GET("/history", s.HistoryPage)
	s.router.POST("/history/search", s.HistorySearch)
	s.router.GET("/", s.Index)
	s.router.GET("/:id", s.Index)
	s.router.GET("/:id/rabbitmq", s.RabbitMQ)
	return s.router.Start(":8080")
}
