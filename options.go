package router

func AllowedMethod(method string) RouteOption {
	return func(r *Route) {
		r.AllowedMethod(method)
	}
}

func Authorized() RouteOption {
	return func(r *Route) {
		r.Authorize()
	}
}
