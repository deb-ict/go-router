package router

func AllowedMethod(method string) RouteOption {
	return func(r *Route) {

	}
}

func AllowedHeader(name string, value string) RouteOption {
	return func(r *Route) {

	}
}

func Authorized() RouteOption {
	return func(r *Route) {

	}
}
