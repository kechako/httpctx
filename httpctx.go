package httpctx

import (
	"net/http"
	"sync"

	"golang.org/x/net/context"
)

var (
	ctxmx sync.Mutex
	ctxs  = make(map[*http.Request]context.Context)
)

// Retrieves a handler that binds the context to the request.
func NewHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxmx.Lock()
		ctxs[r] = context.Background()
		ctxmx.Unlock()

		next.ServeHTTP(w, r)

		ctxmx.Lock()
		delete(ctxs, r)
		ctxmx.Unlock()
	})
}

// Retrieves the context bound to the request.
func Get(r *http.Request) context.Context {
	ctxmx.Lock()
	ctx := ctxs[r]
	ctxmx.Unlock()

	if ctx == nil {
		panic("GetContext passed an unknown http.Request")
	}

	return ctx
}

func Set(r *http.Request, ctx context.Context) {
	ctxmx.Lock()
	if _, ok := ctxs[r]; !ok {
		panic("GetContext passed an unknown http.Request")
	}
	ctxs[r] = ctx
	ctxmx.Unlock()
}
