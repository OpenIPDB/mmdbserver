package mmdbserver

import (
	"errors"
	"net/http"
	"regexp"
)

var (
	reUpdatePath = regexp.MustCompile(`^\/geoip\/databases\/(?P<edition>[^\/]+)\/update$`)
)

type MMDBHandler interface {
	ServeMMDB(*Response, *Request) error
}

type MMDBHandlerFunc func(*Response, *Request) error

func (h MMDBHandlerFunc) ServeMMDB(resp *Response, r *Request) error { return h(resp, r) }

type MMDBUpdateHandler struct {
	HomePage string
	MMDBHandler
}

func (h *MMDBUpdateHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if !reUpdatePath.MatchString(r.URL.Path) {
		http.Redirect(rw, r, h.HomePage, http.StatusTemporaryRedirect)
		return
	}
	response, err := h.invoke(r)
	if err != nil {
		var exception *Error
		if errors.As(err, &exception) {
			http.Error(rw, err.Error(), exception.StatusCode)
		} else {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	_, _ = response.WriteTo(rw)
	return
}

func (h *MMDBUpdateHandler) invoke(r *http.Request) (response *Response, err error) {
	if r.Method != http.MethodGet {
		err = ErrMethodNotAllowed
		return
	}
	request, err := NewRequest(r)
	if err != nil {
		return
	}
	response = NewResponse()
	if err = h.ServeMMDB(response, request); err != nil {
		return
	}
	switch {
	case !response.Valid():
		err = ErrDatabaseNotFound
	case request.EqualHash(response.MD5Hash()):
		err = ErrDatabaseLatest
	}
	return
}
