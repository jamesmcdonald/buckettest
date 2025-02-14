package web

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/jamesmcdonald/buckettest/internal/bucket"
)

//go:embed templates/*
var templateFS embed.FS

type App struct {
	bucket    *bucket.Bucket
	templates *template.Template
	mux       *http.ServeMux
	logger    *slog.Logger
}

func New(bucket *bucket.Bucket, logger *slog.Logger) (*App, error) {
	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return nil, err
	}
	app := &App{
		bucket:    bucket,
		templates: tmpl,
		logger:    logger,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", app.ListObjects)
	mux.HandleFunc("GET /read", app.ReadObject)
	mux.HandleFunc("POST /delete", app.DeleteObject)
	mux.HandleFunc("POST /write", app.WriteObject)
	app.mux = mux

	return app, nil
}

// Web handlers to list, read and delete objects in a bucket
func (a *App) ListObjects(w http.ResponseWriter, r *http.Request) {
	objects, err := a.bucket.ListObjects(r.Context())
	if err != nil {
		http.Error(w, "Error listing objects", http.StatusInternalServerError)
		return
	}
	err = a.templates.ExecuteTemplate(w, "index.html", struct{ Objects []string }{Objects: objects})
}

func (a *App) ReadObject(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")

	object, err := a.bucket.GetObject(r.Context(), path)
	if err != nil {
		http.Error(w, "Error reading object", http.StatusInternalServerError)
		a.logger.ErrorContext(r.Context(), "Error reading object", "err", err, "path", path)
		return
	}
	a.logger.DebugContext(r.Context(), "Object read", "path", path)
	err = a.templates.ExecuteTemplate(w, "read.html", struct{ Path, Content string }{Path: path, Content: string(object)})
	if err != nil {
		a.logger.ErrorContext(r.Context(), "Error rendering template", "err", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func (a *App) DeleteObject(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "No path provided", http.StatusBadRequest)
		a.logger.ErrorContext(r.Context(), "Delete request with no path", "pathvalue", r.PathValue("path"), "formvalue", r.FormValue("path"))
		return
	}
	err := a.bucket.DeleteObject(r.Context(), path)
	if err != nil {
		http.Error(w, "Error deleting object", http.StatusInternalServerError)
		a.logger.ErrorContext(r.Context(), "Error deleting object", "err", err, "path", path)
		return
	}
	a.logger.DebugContext(r.Context(), "Object deleted", "path", path)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *App) WriteObject(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "No path provided", http.StatusBadRequest)
		a.logger.ErrorContext(r.Context(), "Write request with no path", "pathvalue", r.PathValue("path"), "formvalue", r.FormValue("path"))
		return
	}
	content := []byte(r.PostFormValue("content"))
	err := a.bucket.PutObject(r.Context(), path, content)
	if err != nil {
		a.logger.ErrorContext(r.Context(), "Error writing object", "err", err, "path", path)
		http.Error(w, "Error writing object", http.StatusInternalServerError)
		return
	}
	a.logger.DebugContext(r.Context(), "Object written", "path", path)
	http.Redirect(w, r, "/read?path="+path, http.StatusFound)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
