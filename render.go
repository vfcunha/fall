package fall

import (
	"html/template"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
)

var defaultLayout = ""
var partialsTemplates = []string{}

func LoadPartials() {
	err := scanPartials("web/views/partials")
	if err != nil {
		panic(err)
	}
}

func scanPartials(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.html"))
	if err != nil {
		return err
	}
	partialsTemplates = append(partialsTemplates, files...)
	return nil
}

func SetDefaultLayout(layout string) {
	defaultLayout = layout
}

func RenderWithLayout(w http.ResponseWriter, r *http.Request, data any, layout string, templates ...string) {
	patternVal := r.Context().Value(patternContextKey)
	pattern, ok := patternVal.(string)
	if !ok {
		http.Error(w, "Route pattern not found in context", http.StatusInternalServerError)
		return
	}
	tmpl := []string{}
	if templates == nil {
		tmpl = []string{
			"web/views/pages/" + patternToTemplatePath(pattern),
		}
	} else {
		for _, t := range templates {
			tmpl = append(tmpl, "web/views/pages/"+t)
		}
	}
	tmpl = append(tmpl, partialsTemplates...)
	if layout != "" {
		tmpl = append(tmpl, "web/views/layouts/"+layout+".html")
	}

	t, err := template.ParseFiles(tmpl...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if layout == "" {
		err = t.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = t.ExecuteTemplate(w, layout, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Render(w http.ResponseWriter, r *http.Request, data any, tmpl ...string) {
	RenderWithLayout(w, r, data, defaultLayout, tmpl...)
}

var wildcardRegex = regexp.MustCompile(`\{[^}]+\}`)

func patternToTemplatePath(pattern string) string {
	parts := strings.SplitN(pattern, " ", 2) // Separa "METHOD /path"
	if len(parts) != 2 {
		// Se não tiver método (improvável vindo do seu router), ou for só "/", trata como root
		if pattern == "/" {
			return "index.html"
		}
		// Tenta tratar como path direto se não houver método
		cleanPath := strings.Trim(pattern, "/")
		if cleanPath == "" {
			return "index.html"
		}
		// Fallback: assume que é um path direto sem wildcards
		return cleanPath + ".html"
	}

	pathPattern := parts[1] // Pega a parte do path: "/post/{id}"

	// Limpa barras no início/fim para facilitar o split
	cleanPath := strings.Trim(pathPattern, "/")
	if cleanPath == "" { // Raiz "/"
		return "index.html"
	}

	segments := strings.Split(cleanPath, "/")
	lastSegment := segments[len(segments)-1]

	// Verifica se o último segmento é um wildcard
	isLastWildcard := wildcardRegex.MatchString(lastSegment)

	if isLastWildcard {
		// Se for wildcard, remove o último segmento e adiciona "show.html"
		if len(segments) > 1 {
			return filepath.Join(strings.Join(segments[:len(segments)-1], "/"), "show.html")
		}
		// Se era só o wildcard (e.g., "/{id}"), retorna "show.html"
		return "show.html"
	} else {
		if strings.Contains(lastSegment, ".") {
			return cleanPath
		}

		if parts[0] == "POST" {
			return filepath.Join(cleanPath, "/new.html")
		}

		// Senão, assume que é um diretório/coleção, usa index.html
		if len(segments) == 1 {
			return filepath.Join(cleanPath, "index.html")
		} else {
			return cleanPath + ".html"
		}
	}
}
