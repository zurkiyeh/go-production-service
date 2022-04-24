package http

import (
	"context"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// Middleware are intermidiate funcs between request and response. We implement them to use by default to all incoming
// request/responses so that we dont have to repeat ourselves many times
func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next.ServeHTTP(w, r)
	})
}

// Middleware funcs are also very useful when logging info requests.
// This will produce a log that looks like this:
// time="2022-04-24T03:12:25Z" level=info msg="handled request" method=POST path=/api/v1/comment
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(
			log.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
			}).Info("handled request")

		next.ServeHTTP(w, r)
	})
}

// Any function that takes more than 15 seconds will return with the following error:
// error fetching comment by uuid: context deadline exceeded
// The reason this works is that we're passing the context on from the different layers of our applications until we hit the db
// layer with then uses queryRowWithContext which then applies this context

func TimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Here we're calling the context function that will call cancel after the 15 sec window has passed
		// Every request has to return within 15 seconds or else it will timeout with an error
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
