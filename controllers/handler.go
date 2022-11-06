package controllers

import (
	"context"
	"github.com/Arnobkumarsaha/new-server/auth"
	"github.com/go-chi/chi/v5"
	"net/http"
)

/*
If I set alias like , `type ControllerProductResource auth.AuthProductResource`
then, ControllerProductResource will not find the methods (for example AuthMiddleWare() ) of AuthProductResource.
It will only get the fields .
*/

type ControllerProductResource struct {
	*auth.AuthProductResource
}

func (rs *ControllerProductResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(rs.AuthMiddleware)
	r.Get("/", rs.GetAllProducts)
	r.Post("/", rs.CreateProduct)

	r.Route("/{pid}", func(r chi.Router) {
		r.Use(rs.ProductCtx)
		r.Get("/", rs.GetSingleProduct)
		r.Put("/", rs.UpdateProduct)
		r.Delete("/", rs.DeleteProduct)
	})

	return r
}

func (rs *ControllerProductResource) ProductCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// URLParam reads the param (which was set on r.Route()) for url
		// We are setting key-value pair in context to easily get it in the subsequent function calls.
		ctx := context.WithValue(r.Context(), "prod_id", chi.URLParam(r, "pid"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
