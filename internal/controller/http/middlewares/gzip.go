package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

const (
	compression = "gzip"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.Header.Get("Content-Encoding"), compression) {
			reader, err := gzip.NewReader(r.Body)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			r.Body = reader
			defer r.Body.Close()
		}

		if strings.Contains(r.Header.Get("Accept-Encoding"), compression) {
			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			defer gz.Close()
			w.Header().Set("Content-Encoding", compression)
			w = gzipWriter{ResponseWriter: w, Writer: gz}
		}

		next.ServeHTTP(w, r)
	})
}
