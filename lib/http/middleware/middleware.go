package middleware

import (
	"net/http"
)

type Middleware interface {
	Handle(w http.ResponseWriter, r *http.Request, next http.Handler) http.Handler
}

type Collection struct {
	middlewares []Middleware
}

// Инициализация коллекции middleware
func NewCollection() *Collection {
	return &Collection{
		middlewares: make([]Middleware, 0),
	}
}

// Добавление middleware
func (c *Collection) Append(mw Middleware) {
	c.middlewares = append(c.middlewares, mw)
}

func (c *Collection) BuildChain(finalHandler http.Handler) http.Handler {
	if len(c.middlewares) == 0 {
		return finalHandler
	}

	chainedHandler := finalHandler

	for i := len(c.middlewares) - 1; i >= 0; i-- {
		currentMiddleware := c.middlewares[i]
		nextHandler := chainedHandler
		// nil для w и r, т.к. они обрабатываются в Handle
		chainedHandler = currentMiddleware.Handle(nil, nil, nextHandler)
	}

	return chainedHandler
}
