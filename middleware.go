package z

type MiddlewareFunc func(handler HandlerFunc) HandlerFunc
