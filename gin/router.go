package gin

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strings"
	"unsafe"

	"github.com/gin-gonic/gin"
	goplugify "github.com/go-plugify/go-plugify"
)

type HttpContext struct {
	*gin.Context
}

func (ctx *HttpContext) Body() io.ReadCloser {
	body, _ := ctx.Context.Request.GetBody()
	return body
}

type HttpRouter struct {
	engine *gin.Engine
}

func NewHttpRouter(engine *gin.Engine) *HttpRouter {
	return &HttpRouter{engine: engine}
}

func (p *HttpRouter) NewHTTPContext(c context.Context) *HttpContext {
	return &HttpContext{Context: c.(*gin.Context)}
}

func (p *HttpRouter) Add(method, route string, handler func(c goplugify.HttpContext)) {
	p.engine.Handle(strings.ToUpper(method), route, func(c *gin.Context) {
		handler(&HttpContext{Context: c})
	})
}

func (p *HttpRouter) ReplaceHandler(method, path string, handler func(ctx context.Context)) error {
	return ReplaceLastHandler(p.engine, method, path, func(c *gin.Context) {
		handler(c)
	})
}

func (p *HttpRouter) GetHandler(method, path string) (func(ctx context.Context), error) {
	handlers, err := getHandlerSlicePointer(p.engine, method, path)
	if err != nil {
		return nil, err
	}
	handler := (*handlers)[len(*handlers)-1]
	return func(ctx context.Context) {
		handler(ctx.(*gin.Context))
	}, nil
}

func (p *HttpRouter) GetHandlerName(method, path string) (string, error) {
	handlers, err := getHandlerSlicePointer(p.engine, method, path)
	if err != nil {
		return "", err
	}
	if handlers == nil || len(*handlers) == 0 {
		return "", fmt.Errorf("no handlers found for method: %s, route: %s", method, path)
	}
	handler := (*handlers)[len(*handlers)-1]
	handlerName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	return handlerName, nil
}

func (p *HttpRouter) GetRoutes() []gin.RouteInfo {
	return p.engine.Routes()
}

func getHandlerSlicePointer(engine *gin.Engine, method string, route string) (*[]gin.HandlerFunc, error) {
	engineVal := reflect.ValueOf(engine).Elem()
	trees := engineVal.FieldByName("trees")
	if !trees.IsValid() {
		return nil, fmt.Errorf("cannot find route trees")
	}

	for i := range trees.Len() {
		tree := trees.Index(i)
		methodField := getUnexportedField(tree.FieldByName("method")).String()
		if !strings.EqualFold(methodField, method) {
			continue
		}

		root := getUnexportedField(tree.FieldByName("root"))
		handlers := findHandlersInNode(root, route, "")
		if handlers == nil {
			continue
		} else {
			return handlers, nil
		}
	}
	return nil, fmt.Errorf("handler not found for method: %s, route: %s", method, route)
}

func findHandlersInNode(node reflect.Value, target, currentPath string) *[]gin.HandlerFunc {
	if !node.IsValid() {
		return nil
	}

	if node.Kind() == reflect.Ptr {
		if node.IsNil() {
			return nil
		}
		node = node.Elem()
	}
	node = getUnexportedField(node)

	path := getUnexportedField(node.FieldByName("path")).String()
	fullPath := currentPath + path

	if fullPath == target {
		handlersField := getUnexportedField(node.FieldByName("handlers"))
		handlersPtr := (*[]gin.HandlerFunc)(unsafe.Pointer(handlersField.UnsafeAddr()))
		return handlersPtr
	}

	childrenField := getUnexportedField(node.FieldByName("children"))
	for i := range childrenField.Len() {
		child := childrenField.Index(i)
		if child.Kind() == reflect.Ptr && child.IsNil() {
			continue
		}
		if res := findHandlersInNode(child, target, fullPath); res != nil {
			return res
		}
	}
	return nil
}

func ReplaceLastHandler(engine *gin.Engine, method, route string, newHandler gin.HandlerFunc) error {
	handlersPtr, err := getHandlerSlicePointer(engine, method, route)
	if err != nil {
		return fmt.Errorf("failed to get handler slice pointer: %w", err)
	}
	if handlersPtr == nil || len(*handlersPtr) == 0 {
		return fmt.Errorf("no handlers found for method: %s, route: %s", method, route)
	}
	(*handlersPtr)[len(*handlersPtr)-1] = newHandler
	return nil
}

func getUnexportedField(v reflect.Value) reflect.Value {
	if !v.IsValid() {
		return v
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.CanAddr() {
		panic("value is not addressable")
	}
	return reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
}
