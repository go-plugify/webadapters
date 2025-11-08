package echo

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"strings"

	goplugify "github.com/go-plugify/go-plugify"
	"github.com/labstack/echo/v4"
)

type HttpContext struct {
	echoCtx echo.Context
	context.Context
}

func (ctx *HttpContext) GetHeader(key string) string {
	return ctx.echoCtx.Request().Header.Get(key)
}

func (ctx *HttpContext) Body() io.ReadCloser {
	return ctx.echoCtx.Request().Body
}

func (ctx *HttpContext) FormFile(name string) (*multipart.FileHeader, error) {
	_, file, err := ctx.echoCtx.Request().FormFile(name)
	return file, err
}

func (ctx *HttpContext) Query(key string) string {
	return ctx.echoCtx.QueryParam(key)
}

func (ctx *HttpContext) PostForm(key string) string {
	return ctx.echoCtx.FormValue(key)
}

func (ctx *HttpContext) JSON(code int, obj any) {
	ctx.echoCtx.Response().Header().Set("Content-Type", "application/json")
	ctx.echoCtx.Response().WriteHeader(code)
	objData, _ := json.Marshal(obj)
	ctx.echoCtx.Response().Write(objData)
}

type HttpRouter struct {
	echo *echo.Echo
}

func NewHttpRouter(echo *echo.Echo) *HttpRouter {
	return &HttpRouter{echo: echo}
}

func (p *HttpRouter) Add(method, route string, handler func(c goplugify.HttpContext)) {
	p.echo.Add(strings.ToUpper(method), route, func(c echo.Context) error {
		handler(&HttpContext{echoCtx: c, Context: c.Request().Context()})
		return nil
	})
}

func (p *HttpRouter) ReplaceHandler(method, path string, handler func(ctx context.Context)) error {
	panic("not implemented")
}

func (p *HttpRouter) GetHandler(method, path string) (func(ctx context.Context), error) {
	panic("not implemented")
}

func (p *HttpRouter) GetHandlerName(method, path string) (string, error) {
	panic("not implemented")
}
