package http_wrapperfakes

import (
	"net/http"
)

type FakeHttpWrapper struct {
	Url string
	Answer *http.Response
	Err error
}

func (fake *FakeHttpWrapper) Get(r string) (resp *http.Response, err error) {
	fake.Url = r
	return fake.Answer, fake.Err
}