package http_wrapper

import (
	"net/http"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter . HttpWrapperInterface

type HttpWrapperInterface interface {
	Get(request string)  (*http.Response, error)
}

type HttpWrapper struct {
}

func (_ *HttpWrapper) Get(r string) (resp *http.Response, err error) {
	return http.Get(r)
}