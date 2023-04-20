package mweb

type Middleware func(next HandlerFunc) HandlerFunc
