package handler

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist
var staticFiles embed.FS

func spaHandler() http.Handler {
	dist, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(dist))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// tenta servir o arquivo estático
		f, err := dist.Open(r.URL.Path[1:])
		if err != nil {
			// fallback para index.html (SPA routing)
			r.URL.Path = "/"
		} else {
			f.Close()
		}
		fileServer.ServeHTTP(w, r)
	})
}
