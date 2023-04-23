package router

import (
	"net/http"
	"testing"
)

func Test_AllowedMethod(t *testing.T) {
	route := &Route{}
	option := AllowedMethod(http.MethodGet)
	option(route)

	if route.methods == nil || len(route.methods) != 1 || route.methods[0] != http.MethodGet {
		t.Errorf("AllowedMethod() option failed: method not added as allowed to route")
	}
}

func Test_AllowedHeader(t *testing.T) {

}
