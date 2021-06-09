package routes

import (
	"net/http"

	"github.com/go-chi/chi"

	controller "github.com/mktsy/go-login-api/api/controllers"
	db "github.com/mktsy/go-login-api/api/databases"
)

func Router(s *db.Dispatch) http.Handler {
	r := chi.NewRouter()
	r.Route("/users", func(r chi.Router) {
		//CRUD User
		r.Get("/{id}", controller.GetUser(s))
		r.Post("/", controller.CreateUser(s))
		r.Put("/{id}", controller.UpdateUser(s))
		r.Delete("/{id}", controller.DeleteUser(s))
	})
	return r
}
