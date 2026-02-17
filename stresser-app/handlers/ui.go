package handlers

import (
	"net/http"
	"path/filepath"
)

// UIHandler serve a pagina HTML estatica
// GET / → index.html
func UIHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("static", "index.html"))
}

// StaticHandler serve arquivos estaticos (CSS, JS)
// GET /static/* → static/*
func StaticHandler() http.Handler {
	return http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
}