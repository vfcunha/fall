package fall

import (
	"fmt"
	"log/slog"
	"net/http"
	"slices"
)

type Middleware func(http.Handler) http.Handler

type Router struct {
	*http.ServeMux
	prefix string
	chain  []Middleware
}

func NewRouter(prefix string, middlewares ...Middleware) *Router {
	return &Router{
		ServeMux: http.NewServeMux(),
		chain:    middlewares,
	}
}

func (r *Router) Use(mw ...Middleware) {
	r.chain = append(r.chain, mw...)
}

func (r *Router) Group(prefix string, fn func(r *Router)) {
	middlewares := slices.Clone(r.chain)
	fn(&Router{
		prefix:   prefix,
		ServeMux: r.ServeMux,
		chain:    middlewares,
	})
}

func (r *Router) path(path string) string {
	return fmt.Sprintf("/%s%s", r.prefix, path)
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

func (r *Router) Option(path string, fn http.HandlerFunc, mws ...Middleware) {
	r.handle(http.MethodOptions, r.path(path), fn, mws...)
}

func (r *Router) handle(method, path string, fn http.HandlerFunc, mws ...Middleware) {
	slog.Info(fmt.Sprintf("%s %s", method, path))
	r.Handle(fmt.Sprintf("%s %s", method, path), r.wrap(fn, mws...))
}

func (r *Router) wrap(fn http.HandlerFunc, mws ...Middleware) (out http.Handler) {
	out, mwss := http.Handler(fn), append(r.chain, mws...)
	for i := len(mwss) - 1; i >= 0; i-- {
		out = mwss[i](out)
	}
	return out
	// out, mwss := http.Handler(fn), append(slices.Clone(r.chain), mws...)
	// slices.Reverse(mwss)
	// for _, m := range mwss {
	// 	out = m(out)
	// }
	// return
}

func (r *Router) Stack() (out http.Handler) {
	stack := createStack(r.chain...)
	return stack(r)
}

func createStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}
