package router

func AllowedMethod(method string) RouteOption {
	return func(r *Route) {
		r.AllowedMethod(method)
	}
}

func AllowedMethods(method ...string) RouteOption {
	return func(r *Route) {
		r.AllowedMethods(method...)
	}
}

func Authorized(policyName string) RouteOption {
	return func(r *Route) {
		r.Authorize(policyName)
	}
}
