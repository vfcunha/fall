package fall

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
)

type App struct {
	Env    Enviroment `env:"ENV" envDefault:"Development"`
	router *Router
}

func NewApp(env Enviroment, envConfig EnvConfiguration, middlewares ...Middleware) (*App, error) {
	app := App{
		Env:    env,
		router: NewRouter("", middlewares...),
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
	http.ListenAndServe(fmt.Sprintf(":%s", port), a.router.Stack())
	// http.ListenAndServe(fmt.Sprintf(":%s", port), LogRequest(a.router.Stack()))
}

func (a *App) GetRouter() *Router {
	return a.router
}

func LogRequest(handler http.Handler) http.Handler {
	// log.Printf("%s \n", "Log Request")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
	})
}
