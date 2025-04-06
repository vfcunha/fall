package fall

import (
	"fmt"
	"log/slog"
	"net/http"
)

type App struct {
	Env         Enviroment `env:"ENV" envDefault:"Development"`
	router      *Router
	middlewares []Middleware
}

func NewApp(env Enviroment, envConfig EnvConfiguration, middlewares ...Middleware) (*App, error) {
	app := App{
		Env:         env,
		router:      NewRouter(""),
		middlewares: middlewares,
	}
	err := envConfig.Configure(env)
	if err != nil {
		return nil, err
	}

	app.SetControllers(
		ResolveControllers(),
	)

	return &app, nil
}

func (a *App) SetControllers(controllers []Controller) {
	for _, handler := range controllers {
		handler.Configure(a.router)
	}
}

func (a *App) ListenAndServe(port string) {
	slog.Info("HTTP Server started", "listenAddr", port)
	stack := createStack(a.middlewares...)
	http.ListenAndServe(fmt.Sprintf(":%s", port), stack(a.router))
}

func (a *App) GetRouter() *Router {
	return a.router
}

func LogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info(fmt.Sprintf("%s %s %s", r.RemoteAddr, r.Method, r.URL))
		handler.ServeHTTP(w, r)
	})
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
