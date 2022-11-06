package main

import (
	"github.com/Arnobkumarsaha/new-server/auth"
	"github.com/Arnobkumarsaha/new-server/controllers"
	"github.com/Arnobkumarsaha/new-server/schemas"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	schemas.AddRequiredData()
	port := "8080"
	log.Printf("Starting up on http://localhost:%s", port)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/signin", auth.SignIn)
	r.Get("/welcome", auth.Welcome)
	r.Post("/refresh", auth.Refresh)
	r.Post("/logout", auth.LogOut)

	/*
		ControllerProductResource should be initialized if it has a pointer indside.
		But as AuthProductResource has no fields, we hadn't to do that.
		For details :: https://play.golang.org/p/x6q-2HUfTH
	*/
	r.Mount("/products", (&controllers.ControllerProductResource{}).Routes())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) { // Dummy, Default
		w.Write([]byte("Hello World!"))
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
