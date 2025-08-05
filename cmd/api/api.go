package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mafi020/social/internal/store"
)

type dbConfig struct {
	url          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type config struct {
	port string
	db   dbConfig
	env  string
}

type application struct {
	config config
	store  store.Storage
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Patch("/", app.updatePostHandler)
			})
		})
		r.Route("/users", func(r chi.Router) {
			r.Post("/", app.createUserHandler)
			r.Route("/{userID}", func(r chi.Router) {
				r.Get("/", app.getUserHandler)
				r.Delete("/", app.deleteUserHandler)
			})
		})
		r.Route("/comments", func(r chi.Router) {
			r.Post("/", app.createCommentHandler)
			r.Route("/{commentID}", func(r chi.Router) {
				r.Get("/", app.getCommentHandler)
				r.Patch("/", app.updateCommentHandler)
				r.Delete("/", app.deleteCommentHandler)
			})
		})
	})

	return r
}

func (app *application) start(mux http.Handler) error {

	server := &http.Server{
		Addr:         app.config.port,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	log.Printf("Server running on port %s", app.config.port)
	return server.ListenAndServe()
}
