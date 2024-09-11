package web

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type renderParams struct {
	Ctx       echo.Context
	Component templ.Component
	Status    int
}

func render(p renderParams) error {
	stat := http.StatusOK
	if p.Status != 0 {
		stat = p.Status
	}
	c := templ.Handler(p.Component, templ.WithStatus(stat)).Component
	return c.Render(p.Ctx.Request().Context(), p.Ctx.Response())
}
