package resreq

import (
	"github.com/Arnobkumarsaha/new-server/schemas"
	"github.com/go-chi/render"
	"net/http"
)

type ProductResponse struct {
	*schemas.Product
	User *schemas.User `json:"user,omitempty"`
}

var _ render.Renderer = &ProductResponse{}

func (rd *ProductResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("ProductResponse Render() method is called."))
	return nil
}

// We are type casting our Product to render.Renderer interface,
// to call render.Render() or render.RenderList() method with appropriate parameter.

func NewProductResponse(product *schemas.Product) render.Renderer {
	resp := &ProductResponse{Product: product}
	if resp.User == nil {
		if user, _ := getUserByID(resp.OwnerId); user != nil {
			resp.User = user
		}
	}
	return resp
}

func NewProductListResponse(products []*schemas.Product) []render.Renderer {
	var list []render.Renderer
	for _, article := range products {
		list = append(list, NewProductResponse(article))
	}
	return list
}

/*
ABOUT RENDERING (call stack):
render.RenderList & render.Render calls renderer() & Respond().
renderer() calls Render() (not render.render(), but the Render() method of the structure.)

Respond() calls DefaultResponder(), then DefaultResponder() call JSON().
and JSON() actually writes, what We see in the responses.
*/
