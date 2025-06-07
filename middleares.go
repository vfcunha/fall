package fall

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func LogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info(fmt.Sprintf("%s %s %s", r.RemoteAddr, r.Method, r.URL))
		handler.ServeHTTP(w, r)
	})
}

func StaticDir(prefix, publicDir string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, prefix) {
				filePath := filepath.Join(publicDir, strings.TrimPrefix(r.URL.Path, prefix))
				if fileExists(filePath) { // Função para verificar se o arquivo existe
					http.ServeFile(w, r, filePath)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
