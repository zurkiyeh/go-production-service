package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	Router  *mux.Router
	Service CommentService
	Server  *http.Server
}

func NewHandler(service CommentService) *Handler {
	h := &Handler{
		Service: service,
	}
	h.Router = mux.NewRouter()
	h.mapRoutes()
	// the use func is specific to the middleware functionality
	h.Router.Use(JSONMiddleware)
	h.Router.Use(LoggingMiddleware)
	h.Router.Use(TimeoutMiddleware)
	h.Server = &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: h.Router,
	}
	return h
}

func (h *Handler) mapRoutes() {
	// Dont forget the slashes before the api
	h.Router.HandleFunc("/api/v1/comment", h.PostComment).Methods("POST")
	h.Router.HandleFunc("/api/v1/comment/{id}", h.GetComment).Methods("GET")
	h.Router.HandleFunc("/api/v1/comment/{id}", h.UpdateComment).Methods("PUT")
	h.Router.HandleFunc("/api/v1/comment/{id}", h.DeleteComment).Methods("DELETE")
}

func (h *Handler) Serve() error {
	// Graceful shutdown
	// Bacially made this function a non-blocking func so it will spawn in it's own thread and execute until completion
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()

	// Rest of the function will continue to execute regardless of above

	// Create a channel that will receive an os signal
	c := make(chan os.Signal, 1)
	// (BLOCKING ACTION): block until you receive an interrupt signal
	signal.Notify(c, os.Interrupt)

	//(UNBLOCKING ACTION): Read from the channel once you receive something and continue execution
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// When 15 seconds have elapsed, the cancel func will be called
	defer cancel()
	// Shutsdown gracefully: finishes up all requests and stops receving new requests
	h.Server.Shutdown(ctx)
	log.Println("shut down gracefully")
	return nil

}
