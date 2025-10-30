package gorilla

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"

	goplugify "github.com/go-plugify/go-plugify"
	"github.com/gorilla/mux"
)

type HttpContext struct {
	w http.ResponseWriter
	r *http.Request
}

func (ctx *HttpContext) GetHeader(key string) string {
	return ctx.r.Header.Get(key)
}

func (ctx *HttpContext) Body() io.ReadCloser {
	return ctx.r.Body
}

func (ctx *HttpContext) FormFile(name string) (*multipart.FileHeader, error) {
	_, file, err := ctx.r.FormFile(name)
	return file, err
}

func (ctx *HttpContext) Query(key string) string {
	return ctx.r.URL.Query().Get(key)
}

func (ctx *HttpContext) JSON(code int, obj any) {
	ctx.w.Header().Set("Content-Type", "application/json")
	ctx.w.WriteHeader(code)
	objData, _ := json.Marshal(obj)
	ctx.w.Write(objData)
}

func (ctx *HttpContext) PostForm(key string) string {
	return ctx.r.PostFormValue(key)
}

type HttpRouter struct {
	app *mux.Router
}

func NewHttpRouter(app *mux.Router) *HttpRouter {
	return &HttpRouter{app: app}
}

func (p *HttpRouter) Add(route string, handler func(c goplugify.HttpContext)) {
	p.app.Handle(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(&HttpContext{w: w, r: r})
	})).Methods("POST")
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
