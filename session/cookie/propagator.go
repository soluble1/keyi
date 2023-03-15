package cookie

import (
	"net/http"
)

type CookieOption func(cookie *http.Cookie)

type PropagatorOption func(p *Propagator)

type Propagator struct {
	cookieName string
	cookieOpt  CookieOption
}

func WithCookieOpt(opt CookieOption) PropagatorOption {
	return func(p *Propagator) {
		p.cookieOpt = opt
	}
}

func NewPropagator(cookieName string, opts ...PropagatorOption) *Propagator {
	pro := &Propagator{
		cookieName: cookieName,
		cookieOpt: func(cookie *http.Cookie) {

		},
	}

	for _, opt := range opts {
		opt(pro)
	}

	return pro
}

func (p *Propagator) Inject(id string, writer http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:  p.cookieName,
		Value: id,
	}
	p.cookieOpt(cookie)
	http.SetCookie(writer, cookie)
	return nil
}

func (p *Propagator) Extract(req *http.Request) (string, error) {
	cookie, err := req.Cookie(p.cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (p *Propagator) Remove(writer http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:   p.cookieName,
		MaxAge: -1,
	}
	p.cookieOpt(cookie)
	http.SetCookie(writer, cookie)
	return nil
}
