package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/Arnobkumarsaha/new-server/errorhandler"
	"github.com/Arnobkumarsaha/new-server/schemas"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

func (rs *ControllerProductResource) ParseProductFromRequestBody(w http.ResponseWriter, r *http.Request) schemas.Product {
	var newProduct schemas.Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		_ = render.Render(w, r, errorhandler.ErrInvalidRequest(err))
	}
	if newProduct.ID == nil { // ID not given
		err = fmt.Errorf("ID field is needed")
		_ = render.Render(w, r, errorhandler.ErrInvalidRequest(err))
	}
	return newProduct
}

func setDefaultHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func getIDFromRequestContext(r *http.Request) string {
	return r.Context().Value("prod_id").(string)
}

func isEqual(a interface{}, b interface{}) bool {
	var A, B []byte

	switch u := a.(type) {
	case int:
		A = []byte(strconv.Itoa(u))
	case string:
		A = []byte(u)
	case int64:
		A = []byte(strconv.FormatInt(u, 10))
	case *int64:
		A = []byte(strconv.FormatInt(*u, 10))
	}

	switch u := b.(type) {
	case int:
		B = []byte(strconv.Itoa(u))
	case string:
		B = []byte(u)
	case int64:
		B = []byte(strconv.FormatInt(u, 10))
	case *int64:
		B = []byte(strconv.FormatInt(*u, 10))
	}

	for len(A) != len(B) {
		return false
	}
	for i, j := range A {
		if B[i] != j {
			return false
		}
	}
	return true
}
