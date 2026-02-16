package handlers

import (
	"net/http"
	"path/filepath"
)

// UIHandler serve a página HTML estática
// GET / → index.html
func UIHandler(w http.ResponseWriter, r *http.Request) {
	// Só serve o index.html para a rota raiz
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, filepath.Join("static", "index.html"))
}