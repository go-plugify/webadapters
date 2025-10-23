package iris

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"

	goplugify "github.com/go-plugify/go-plugify"
	"github.com/kataras/iris/v12"
)

type HttpContext struct {
	iris.Context
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
	return ctx.URLParam(key)
}

func (ctx *HttpContext) JSON(code int, obj any) {
	ctx.ResponseWriter().Header().Set("Content-Type", "application/json")
	ctx.StatusCode(code)
	objData, _ := json.Marshal(obj)
	ctx.Write(objData)
}

type HttpRouter struct {
	app *iris.Application
}

func NewHttpRouter(app *iris.Application) *HttpRouter {
	return &HttpRouter{app: app}
}

func (p *HttpRouter) Add(route string, handler func(c goplugify.HttpContext)) {
	p.app.Post(route, func(ctx iris.Context) {
		handler(&HttpContext{Context: ctx})
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
