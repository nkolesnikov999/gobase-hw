package auth

import (
	"api/orders/configs"
	"api/orders/internal/user"
	"api/orders/pkg/jwt"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
)

type AuthService struct {
	UserRepository *user.UserRepository
	Config         *configs.Config
}

func NewAuthService(userRepository *user.UserRepository, cfg *configs.Config) *AuthService {
	return &AuthService{UserRepository: userRepository, Config: cfg}
}

func (service *AuthService) generateSessionId() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (service *AuthService) generateCode() (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}

// Auth validates phone, ensures user exists, sets session and code, logs code, returns sessionId
func (service *AuthService) Auth(phone string) (string, error) {
	if len(phone) == 0 {
		return "", errors.New(ErrInvalidPhone)
	}
	// Find or create user by phone
	existedUser, err := service.UserRepository.FindByPhone(phone)
	if err != nil || existedUser == nil {
		// create new user
		created, createErr := service.UserRepository.Create(&user.User{Phone: phone})
		if createErr != nil {
			return "", createErr
		}
		existedUser = created
	}

	sessionId, err := service.generateSessionId()
	if err != nil {
		return "", err
	}
	code, err := service.generateCode()
	if err != nil {
		return "", err
	}

	existedUser.SessionId = sessionId
	existedUser.Code = code
	if saveErr := service.UserRepository.Save(existedUser); saveErr != nil {
		return "", saveErr
	}

	fmt.Printf("Auth code for %s (session %s): %04d\n", phone, sessionId, code)
	return sessionId, nil
}

// Verify checks code by sessionId, clears session/code if invalid, issues JWT on success
func (service *AuthService) Verify(sessionId string, code int) (string, error) {
	existedUser, err := service.UserRepository.FindBySessionId(sessionId)
	if err != nil || existedUser == nil {
		return "", errors.New(ErrSessionNotFound)
	}
	if existedUser.Code != code {
		// clear session and code on invalid attempt
		existedUser.SessionId = ""
		existedUser.Code = 0
		_ = service.UserRepository.Save(existedUser)
		return "", errors.New(ErrInvalidCode)
	}
	// success: issue JWT and clear session/code
	tok, err := jwt.NewJWT(service.Config.Auth.Secret).Create(existedUser.Phone)
	if err != nil {
		return "", err
	}
	existedUser.SessionId = ""
	existedUser.Code = 0
	_ = service.UserRepository.Save(existedUser)
	return tok, nil
}
