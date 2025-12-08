package routes_middleware

import "net/http"

type ExampleMiddleware struct{}

func NewExampleMiddleware() *ExampleMiddleware {
	return &ExampleMiddleware{}
}

func (mw *ExampleMiddleware) Handle(w http.ResponseWriter, r *http.Request, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логика до
		next.ServeHTTP(w, r)
		// Логика после
	})
}
