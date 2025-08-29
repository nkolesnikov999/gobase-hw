package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"sync"

	"github.com/jordan-wright/email"

	"adv/verify/configs"
)

type VerifyHandler struct {
	config      configs.Config
	mu          sync.Mutex
	hashToEmail map[string]string
}

func NewVerifyHandler(conf configs.Config) *VerifyHandler {
	return &VerifyHandler{
		config:      conf,
		hashToEmail: make(map[string]string),
	}
}

type sendRequest struct {
	Email string `json:"email"`
}

func (h *VerifyHandler) Send(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req sendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Email) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	hash, err := generateHash()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not generate hash"})
		return
	}

	h.mu.Lock()
	h.hashToEmail[hash] = req.Email
	h.mu.Unlock()

	linkHost := r.Host
	if linkHost == "" {
		linkHost = "localhost:8081"
	}
	verifyURL := fmt.Sprintf("http://%s/verify/%s", linkHost, hash)

	e := email.NewEmail()
	e.From = h.config.SMTP.Email
	e.To = []string{req.Email}
	e.Subject = "Email verification"
	e.Text = []byte(fmt.Sprintf("Please verify your email by visiting: %s", verifyURL))

	host := h.config.SMTP.Address
	if idx := strings.Index(host, ":"); idx > 0 {
		host = host[:idx]
	}
	var auth smtp.Auth
	if h.config.SMTP.Email != "" || h.config.SMTP.Password != "" {
		auth = smtp.PlainAuth("", h.config.SMTP.Email, h.config.SMTP.Password, host)
	}

	if err := e.Send(h.config.SMTP.Address, auth); err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "failed to send email"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "sent"})
}

func (h *VerifyHandler) Verify(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	const prefix = "/verify/"
	if !strings.HasPrefix(path, prefix) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	hash := strings.TrimPrefix(path, prefix)
	if hash == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing hash"})
		return
	}

	h.mu.Lock()
	emailAddr, ok := h.hashToEmail[hash]
	if ok {
		delete(h.hashToEmail, hash)
	}
	h.mu.Unlock()

	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "invalid or expired hash"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "verified", "email": emailAddr})
}

func generateHash() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
