package web_copy

type Middleware func(next HandlerFunc) HandlerFunc
