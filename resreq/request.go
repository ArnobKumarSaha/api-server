package resreq

import (
	"errors"
	"fmt"
	"github.com/Arnobkumarsaha/new-server/schemas"
	"net/http"
	"strings"
)

type ProductRequest struct {
	*schemas.Product
	User *schemas.User `json:"user,omitempty"`
}

func (a *ProductRequest) Bind(r *http.Request) error {
	fmt.Println("ArticleRequest Bind() method is called.")
	if a.Product == nil {
		return errors.New("missing required Article fields")
	}

	a.Product.Title = strings.ToLower(a.Product.Title)
	return nil
}

func getUserByID(ownerID *int64) (*schemas.User, error) {
	if ownerID == nil {
		return nil, errors.New("ownerID is not set in the product")
	}
	for _, u := range schemas.Users {
		if u.ID == ownerID {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}
