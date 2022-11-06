package controllers

import (
	"errors"
	"fmt"
	"github.com/Arnobkumarsaha/new-server/errorhandler"
	"github.com/Arnobkumarsaha/new-server/resreq"
	"github.com/Arnobkumarsaha/new-server/schemas"
	"github.com/go-chi/render"
	"net/http"
)

func (rs *ControllerProductResource) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	setDefaultHeader(w)
	prods := resreq.NewProductListResponse(schemas.Products)
	if err := render.RenderList(w, r, prods); err != nil {
		_ = render.Render(w, r, errorhandler.ErrRender(err))
		return
	}
}

func (rs *ControllerProductResource) GetSingleProduct(w http.ResponseWriter, r *http.Request) {
	setDefaultHeader(w)
	id := getIDFromRequestContext(r)

	for _, p := range schemas.Products {
		if isEqual(id, p.ID) {
			prod := resreq.NewProductResponse(p)
			if err := render.Render(w, r, prod); err != nil {
				_ = render.Render(w, r, errorhandler.ErrRender(err))
				return
			}
			return
		}
	}
	err := errors.New(fmt.Sprintf("no products found with ID = %v", id))
	_ = render.Render(w, r, errorhandler.ErrNotFound(err))
}

func (rs *ControllerProductResource) CreateProduct(w http.ResponseWriter, r *http.Request) {
	setDefaultHeader(w)
	newProduct := rs.ParseProductFromRequestBody(w, r)
	schemas.Products = append(schemas.Products, &newProduct)
	errorhandler.Write(w, fmt.Sprintf("%v added !", newProduct))
}

func (rs *ControllerProductResource) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	setDefaultHeader(w)
	updatedProduct := rs.ParseProductFromRequestBody(w, r)
	id := getIDFromRequestContext(r)

	if !isEqual(updatedProduct.ID, id) {
		err := errors.New(fmt.Sprintf("ID can not be changed when updating."))
		_ = render.Render(w, r, errorhandler.ErrInvalidRequest(err))
		return
	}

	for idx, p := range schemas.Products {
		if isEqual(id, p.ID) {
			schemas.Products[idx] = &updatedProduct
			errorhandler.Write(w, "product updated")
			return
		}
	}

	err := errors.New(fmt.Sprintf("no products found with ID = %v", id))
	_ = render.Render(w, r, errorhandler.ErrNotFound(err))
}

func (rs *ControllerProductResource) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	setDefaultHeader(w)
	id := getIDFromRequestContext(r)

	for idx, p := range schemas.Products {
		if isEqual(id, p.ID) {
			schemas.Products = append(schemas.Products[:idx], schemas.Products[idx+1:]...)
			errorhandler.Write(w, "product deleted")
			return
		}
	}
	err := errors.New(fmt.Sprintf("no products found with ID = %v", id))
	_ = render.Render(w, r, errorhandler.ErrNotFound(err))
}
