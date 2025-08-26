package api

import (
	"net/http"

	"adv/verify/configs"
)

// ApiHandler wires all API routes.
type ApiHandler struct {
	Config configs.Config
}

// NewApiHandler registers all HTTP handlers on the provided mux.
func NewApiHandler(mux *http.ServeMux, h ApiHandler) {
	verifyHandler := NewVerifyHandler(h.Config)

	mux.HandleFunc("POST /send", verifyHandler.Send)
	mux.HandleFunc("GET /verify/", verifyHandler.Verify)
}
