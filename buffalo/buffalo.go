package buffalo

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"strings"

	goplugify "github.com/go-plugify/go-plugify"
	"github.com/gobuffalo/buffalo"
)

type HttpContext struct {
	buffalo.Context
}

func (ctx *HttpContext) GetHeader(key string) string {
	return ctx.Request().Header.Get(key)
}

func (ctx *HttpContext) Body() io.ReadCloser {
	return ctx.Request().Body
}

func (ctx *HttpContext) FormFile(name string) (*multipart.FileHeader, error) {
	_, file, err := ctx.Request().FormFile(name)
	return file, err
}

func (ctx *HttpContext) Query(key string) string {
	return ctx.Request().URL.Query().Get(key)
}

func (ctx *HttpContext) PostForm(key string) string {
	return ctx.Request().PostFormValue(key)
}

func (ctx *HttpContext) JSON(code int, obj any) {
	ctx.Response().Header().Set("Content-Type", "application/json")
	ctx.Response().WriteHeader(code)
	objData, _ := json.Marshal(obj)
	ctx.Response().Write(objData)
}

type HttpRouter struct {
	app *buffalo.App
}

func NewHttpRouter(app *buffalo.App) *HttpRouter {
	return &HttpRouter{app: app}
}

func (p *HttpRouter) Add(method, route string, handler goplugify.Handler) {
	handlerFunc := func(c buffalo.Context) error {
		handler(&HttpContext{Context: c})
		return nil
	}
	switch strings.ToLower(method) {
	case "get":
		p.app.GET(route, handlerFunc)
	case "post":
		p.app.POST(route, handlerFunc)
	case "put":
		p.app.PUT(route, handlerFunc)
	case "delete":
		p.app.DELETE(route, handlerFunc)
	case "patch":
		p.app.PATCH(route, handlerFunc)
	case "head":
		p.app.HEAD(route, handlerFunc)
	case "options":
		p.app.OPTIONS(route, handlerFunc)
	}
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
