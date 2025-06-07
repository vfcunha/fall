package fall

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"
)

type Middleware func(http.Handler) http.Handler

type Router struct {
	*http.ServeMux
	prefix string
	chain  []Middleware
}

func NewRouter(prefix string, middlewares ...Middleware) *Router {
	return &Router{
		prefix:   prefix,
		ServeMux: http.NewServeMux(),
		chain:    middlewares,
	}
}

type contextKey string

const patternContextKey contextKey = "fall.pattern"

func (r *Router) Use(mw ...Middleware) {
	r.chain = append(r.chain, mw...)
}

func (r *Router) Group(prefix string, fn func(r *Router)) {
	middlewares := slices.Clone(r.chain)
	fn(&Router{
		prefix:   r.path(prefix),
		ServeMux: r.ServeMux,
		chain:    middlewares,
	})
}

func (r *Router) path(path string) string {
	return fmt.Sprintf("%s%s", r.format(r.prefix), r.format(path))
}

func (r *Router) format(path string) string {
	if path == "" {
		return ""
	}
	if strings.HasPrefix(path, "/") {
		return path
	}
	return "/" + path
}

func (r *Router) Get(path string, fn http.HandlerFunc, mws ...Middleware) {
	r.handle(http.MethodGet, r.path(path), fn, mws...)
}

func (r *Router) Post(path string, fn http.HandlerFunc, mws ...Middleware) {
	r.handle(http.MethodPost, r.path(path), fn, mws...)
}

func (r *Router) Put(path string, fn http.HandlerFunc, mws ...Middleware) {
	r.handle(http.MethodPut, r.path(path), fn, mws...)
}

func (r *Router) Delete(path string, fn http.HandlerFunc, mws ...Middleware) {
	r.handle(http.MethodDelete, r.path(path), fn, mws...)
}

func (r *Router) Patch(path string, fn http.HandlerFunc, mws ...Middleware) {
	r.handle(http.MethodPatch, r.path(path), fn, mws...)
}

func (r *Router) Options(path string, fn http.HandlerFunc, mws ...Middleware) {
	r.handle(http.MethodOptions, r.path(path), fn, mws...)
}

func (r *Router) handle(method, path string, fn http.HandlerFunc, mws ...Middleware) {
	fullPattern := fmt.Sprintf("%s %s", method, path)
	slog.Info(fullPattern)
	r.Handle(fmt.Sprintf("%s %s", method, path), r.wrap(fn, fullPattern, mws...))
}

func (r *Router) wrap(fn http.HandlerFunc, routePattern string, mws ...Middleware) (out http.Handler) {
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), patternContextKey, routePattern)
		req = req.WithContext(ctx)
		fn(w, req)
	})

	out, mwss := http.Handler(baseHandler), append(r.chain, mws...)
	for i := len(mwss) - 1; i >= 0; i-- {
		out = mwss[i](out)
	}
	return out
}
