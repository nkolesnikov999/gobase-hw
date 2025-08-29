package auth

type AuthRequest struct {
	Phone string `json:"phone" validate:"required,phone"`
}

type AuthResponse struct {
	SessionId string `json:"sessionId"`
}

type VerifyRequest struct {
	SessionId string `json:"sessionId" validate:"required"`
	Code      int    `json:"code" validate:"required"`
}

type VerifyResponse struct {
	Token string `json:"token"`
}
