package router

import (
	"net/http"
)

type RouteOption func(*Route)

type Route struct {
	Handler http.Handler
}
