package auth

import (
	"api/orders/configs"
	"api/orders/pkg/req"
	"api/orders/pkg/res"
	"net/http"
)

type AuthHandlerDeps struct {
	*configs.Config
	Service *AuthService
}

type AuthHandler struct {
	*configs.Config
	AuthService *AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.Service,
	}
	router.HandleFunc("POST /auth", handler.Auth())
	router.HandleFunc("POST /auth/verify", handler.Verify())
}

func (handler *AuthHandler) Auth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[AuthRequest](&w, r)
		if err != nil {
			return
		}
		// Validation already performed in HandleBody via validate tag. Proceed.
		sessionId, err := handler.AuthService.Auth(body.Phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data := AuthResponse{SessionId: sessionId}
		res.Json(w, data, http.StatusOK)
	}
}

func (handler *AuthHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[VerifyRequest](&w, r)
		if err != nil {
			return
		}
		token, err := handler.AuthService.Verify(body.SessionId, body.Code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		data := VerifyResponse{Token: token}
		res.Json(w, data, http.StatusOK)
	}
}
