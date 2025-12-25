package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kararnab/authdemo/pkg/iam/policy"
)

func NewRouter(
	auth *Handlers,
	books *BookHandlers,
	keyRotationHandler *KeyRotationHandler,
	metricsHandler http.Handler,
) http.Handler {
	r := chi.NewRouter()

	// ================================
	// Metrics (public)
	// ================================
	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metricsHandler.ServeHTTP(w, r)
	})

	r.Route("/api", func(r chi.Router) {

		// Public
		r.Post("/register", auth.Register)
		r.Post("/login", auth.Login)
		r.Post("/refresh", auth.Refresh)
		r.Post("/logout", auth.Logout)

		// Protected
		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(auth.IAM))

			r.Route("/books", func(r chi.Router) {
				r.Get("/", books.List)          // GET /api/books
				r.Post("/", books.Create)       // POST /api/books
				r.Get("/{id}", books.Get)       // GET /api/books/{id}
				r.Put("/{id}", books.Update)    // PUT /api/books/{id}
				r.Delete("/{id}", books.Delete) // DELETE /api/books/{id}
			})
		})
	})

	r.Route("/admin", func(r chi.Router) {
		r.Use(AuthMiddleware(auth.IAM))
		r.Use(PolicyMiddleware(
			auth.IAM,
			policy.Action(policy.Admin),
			policy.ResourceContext{Type: policy.Admin},
		))

		r.Post("/keys/rotate", keyRotationHandler.Rotate)
	})

	return r
}
