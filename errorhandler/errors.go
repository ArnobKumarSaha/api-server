package errorhandler

import (
	"fmt"
	"github.com/go-chi/render"
	"net/http"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime errorhandler
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`                 // user-level status message
	AppCode    int64  `json:"code,omitempty"`         // application-specific errorhandler code
	ErrorText  string `json:"errorhandler,omitempty"` // application-level errorhandler message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

func ErrNotFound(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 404,
		StatusText:     "Resource not found.",
		ErrorText:      err.Error(),
	}
}

func Write(w http.ResponseWriter, s string) {
	_, err := w.Write([]byte(s))
	if err != nil {
		_ = fmt.Errorf("error occured when writing %v", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
