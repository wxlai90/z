package z

type MiddlewareFunc func(next HandlerFunc) HandlerFunc
