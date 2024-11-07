package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Logger struct {
	handler http.Handler
}

func NewLogger(handlerToWrap http.Handler) *Logger {
	return &Logger{handlerToWrap}
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := context.WithValue(context.Background(), "reqID", "1234")
	l.handler.ServeHTTP(w, r.WithContext(ctx))
	log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
}
