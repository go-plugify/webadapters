package fasthttp

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"

	"github.com/buaazp/fasthttprouter"
	goplugify "github.com/go-plugify/go-plugify"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type HttpContext struct {
	*fasthttp.RequestCtx
}

func (ctx *HttpContext) GetHeader(key string) string {
	return string(ctx.Request.Header.Peek(key))
}

func (ctx *HttpContext) Body() io.ReadCloser {
	return bodyReaderCloser{Reader: ctx.Request.BodyStream(), req: &ctx.Request}
}

type bodyReaderCloser struct {
	io.Reader
	req *fasthttp.Request
}

func (r bodyReaderCloser) Close() error {
	return r.req.CloseBodyStream()
}

func (ctx *HttpContext) FormFile(name string) (*multipart.FileHeader, error) {
	formReader, err := ctx.Request.MultipartForm()
	if err != nil {
		return nil, err
	}
	files := formReader.File[name]
	if len(files) == 0 {
		return nil, fiber.ErrUnprocessableEntity
	}
	file := files[0]
	return file, err
}

func (ctx *HttpContext) Query(key string) string {
	return ctx.Query(key)
}

func (ctx *HttpContext) JSON(code int, obj any) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.SetStatusCode(code)
	objData, _ := json.Marshal(obj)
	ctx.Write(objData)
}

func (ctx *HttpContext) PostForm(key string) string {
	return string(ctx.FormValue(key))
}

type HttpRouter struct {
	app *fasthttprouter.Router
}

func NewHttpRouter(app *fasthttprouter.Router) *HttpRouter {
	return &HttpRouter{app: app}
}

func (p *HttpRouter) Add(route string, handler func(c goplugify.HttpContext)) {
	p.app.POST(route, func(ctx *fasthttp.RequestCtx) {
		handler(&HttpContext{RequestCtx: ctx})
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
