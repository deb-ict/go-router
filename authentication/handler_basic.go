package authentication

import (
	"encoding/base64"
	"net/http"
	"strings"
)

type BasicAuthenticationHandlerOption func(*BasicAuthenticationHandler)

type BasicAuthenticationHandler struct {
}

func (h *BasicAuthenticationHandler) HandleAuthentication(r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return
	}
	if len(auth) < len(BasicTokenPrex) || !equalFold(auth[:len(BasicTokenPrex)], BasicTokenPrex) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(BasicTokenPrex):])
	if err != nil {
		return
	}
	cs := string(c)
	username, password, ok := strings.Cut(cs, ":")
	if !ok {
		return
	}

	if username != "my-user" || password != "my-pass" {
		return
	}
}
