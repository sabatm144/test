package routes

import (
	"context"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	controller "test_1/server/controllers"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		t1 := time.Now()

		next.ServeHTTP(w, r)

		t2 := time.Now()
		log.Printf(" [%s] %s %s", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				log.Print(string(debug.Stack()))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func wrapHandler(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func noDirListingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]
		f, err := os.Stat(path)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if f != nil && f.IsDir() {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func HTTPRouteConfig() *httprouter.Router {

	router := httprouter.New()
	handler := alice.New(loggingHandler, recoverHandler)
	router.GET("/readxlsx/:key", wrapHandler(handler.ThenFunc(controller.ReadXLSX)))

	return router
}
