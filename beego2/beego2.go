package gf2

import (
	"context"
	"io"
	"mime/multipart"
	"strings"

	"github.com/beego/beego/v2/server/web"
	beegoCtx "github.com/beego/beego/v2/server/web/context"
	goplugify "github.com/go-plugify/go-plugify"
)

type HttpContext struct {
	*beegoCtx.Context
}

func (ctx *HttpContext) GetHeader(key string) string {
	return ctx.Request.Header.Get(key)
}

func (ctx *HttpContext) Body() io.ReadCloser {
	return ctx.Request.Body
}

func (ctx *HttpContext) FormFile(name string) (*multipart.FileHeader, error) {
	_, file, err := ctx.Request.FormFile(name)
	return file, err
}

func (ctx *HttpContext) Query(key string) string {
	return ctx.Request.URL.Query().Get(key)
}

func (ctx *HttpContext) JSON(code int, obj any) {
	ctx.JSON(code, obj)
}

func (ctx *HttpContext) PostForm(key string) string {
	return ctx.Request.PostFormValue(key)
}

type HttpRouter struct {
	app *web.HttpServer
}

func NewHttpRouter(app *web.HttpServer) *HttpRouter {
	return &HttpRouter{app: app}
}

func (p *HttpRouter) Add(method, route string, handler func(c goplugify.HttpContext)) {
	handlerFunc := func(ctx *beegoCtx.Context) {
		handler(&HttpContext{Context: ctx})
	}
	switch strings.ToLower(method) {
	case "get":
		p.app.Get(route, handlerFunc)
	case "post":
		p.app.Post(route, handlerFunc)
	case "put":
		p.app.Put(route, handlerFunc)
	case "delete":
		p.app.Delete(route, handlerFunc)
	case "patch":
		p.app.Patch(route, handlerFunc)
	case "head":
		p.app.Head(route, handlerFunc)
	case "options":
		p.app.Options(route, handlerFunc)
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
