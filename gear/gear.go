package gear

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"

	goplugify "github.com/go-plugify/go-plugify"
	"github.com/teambition/gear"
)

type HttpContext struct {
	*gear.Context
}

func (ctx *HttpContext) GetHeader(key string) string {
	return ctx.Req.Header.Get(key)
}

func (ctx *HttpContext) Body() io.ReadCloser {
	return ctx.Req.Body
}

func (ctx *HttpContext) FormFile(name string) (*multipart.FileHeader, error) {
	_, file, err := ctx.Req.FormFile(name)
	return file, err
}

func (ctx *HttpContext) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

func (ctx *HttpContext) JSON(code int, obj any) {
	ctx.Res.Header().Set("Content-Type", "application/json")
	ctx.Res.WriteHeader(code)
	objData, _ := json.Marshal(obj)
	ctx.Res.Write(objData)
}

func (ctx *HttpContext) PostForm(key string) string {
	return ctx.Req.PostFormValue(key)
}

type HttpRouter struct {
	router *gear.Router
}

func NewHttpRouter(router *gear.Router) *HttpRouter {
	return &HttpRouter{router: router}
}

func (p *HttpRouter) Add(route string, handler func(c goplugify.HttpContext)) {
	p.router.Post(route, func(ctx *gear.Context) error {
		handler(&HttpContext{Context: ctx})
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
