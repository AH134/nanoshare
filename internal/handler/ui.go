package handler

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

func SpaHandler(embedded fs.FS) (http.Handler, error) {
	distFS, err := fs.Sub(embedded, "dist")
	if err != nil {
		return nil, err
	}

	fileServer := http.FileServerFS(distFS)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cleanPath := path.Clean(r.URL.Path)
		filePath := strings.TrimPrefix(cleanPath, "/")

		if filePath == "" {
			filePath = "index.html"
		}

		fileInfo, err := fs.Stat(distFS, filePath)
		if err != nil || fileInfo.IsDir() {
			r.URL.Path = "/"
		}

		fileServer.ServeHTTP(w, r)
	}), nil
}
