package main

import (
	"log"
	"net/http"

	gorillaCtx "github.com/gorilla/context"
	googleCtx "golang.org/x/net/context"
)

// Define a key using a custom int type to avoid name collisions
type key int

const ctxKey key = 0

// middlewareOne creates a google context
func middlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gorillaCtx.Set(r, ctxKey, googleCtx.Background())
		next.ServeHTTP(w, r)
	})
}

// middlewareTwo adds 1st token
func middlewareTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := gorillaCtx.Get(r, ctxKey).(googleCtx.Context)
		newCtx := googleCtx.WithValue(ctx, "tokenKey1", "tokenValue1")
		// need to put it back!!
		gorillaCtx.Set(r, ctxKey, newCtx)
		next.ServeHTTP(w, r)
	})
}

// middlewareThree adds 2nd token
func middlewareThree(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := gorillaCtx.Get(r, ctxKey).(googleCtx.Context)
		newCtx := googleCtx.WithValue(ctx, "tokenKey2", "tokenValue2")
		// need to put it back!!
		gorillaCtx.Set(r, ctxKey, newCtx)
		next.ServeHTTP(w, r)
	})
}

func final(w http.ResponseWriter, r *http.Request) {
	ctx := gorillaCtx.Get(r, ctxKey).(googleCtx.Context)
	t1 := ctx.(googleCtx.Context).Value("tokenKey1").(string)
	t2 := ctx.(googleCtx.Context).Value("tokenKey2").(string)
	log.Println("In the final handler", "tokenKey1", t1, "tokenKey2", t2)
	w.Write([]byte("OK"))
}

func main() {
	finalHandler := http.HandlerFunc(final)

	http.Handle("/", middlewareOne(middlewareTwo(middlewareThree(finalHandler))))
	log.Fatal(http.ListenAndServe(":3000", nil))
}
