package gf2

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"strings"

	goplugify "github.com/go-plugify/go-plugify"
	"github.com/gogf/gf/v2/net/ghttp"
)

type HttpContext struct {
	*ghttp.Request
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
	ctx.Response.Header().Set("Content-Type", "application/json")
	ctx.Response.WriteHeader(code)
	objData, _ := json.Marshal(obj)
	ctx.Response.Write(objData)
}

func (ctx *HttpContext) PostForm(key string) string {
	return ctx.Request.PostFormValue(key)
}

type HttpRouter struct {
	app *ghttp.Server
}

func NewHttpRouter(app *ghttp.Server) *HttpRouter {
	return &HttpRouter{app: app}
}

func (p *HttpRouter) Add(method, route string, handler func(c goplugify.HttpContext)) {
	p.app.BindHandler(strings.ToUpper(method)+":"+route, func(r *ghttp.Request) {
		handler(&HttpContext{Request: r})
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
