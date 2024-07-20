package http_server

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
)

type middleware func(next http.Handler) http.Handler

func withContentTypeApplicationJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func withRequestBodyLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := bodySizeLimiter(1024/2, r.Body)
		if err != nil {
			badRequest(w, "client", err.Error())
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(w, r)
	})
}

type ResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (r *ResponseWriter) Write(p []byte) (int, error) {
	r.buf.Write(p)
	return r.ResponseWriter.Write(p)
}

func withLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(reqBody))
		mrw := &ResponseWriter{w, &bytes.Buffer{}}
		next.ServeHTTP(mrw, r)
		log.Printf("Info: %s %s %s \nReq Body: %v \nRes Body: %v", r.RemoteAddr, r.Method, r.URL.Path, string(reqBody), mrw.buf.String())
	})
}

func registerMiddlewares(next http.Handler) http.Handler {
	var middlewares = []middleware{
		withContentTypeApplicationJson,
		withRequestBodyLimit,
		withLogger,
	}
	for _, middleware := range middlewares {
		next = middleware(next)
	}
	return next
}

func bodySizeLimiter(size int, reqBody io.Reader) (body []byte, err error) {
	body = make([]byte, size)
	reader := io.LimitReader(reqBody, int64(size))
	n, _ := reader.Read(body)
	if n == size {
		return nil, errors.New("request body is too long")
	}
	return body, nil
}
