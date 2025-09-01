package jwt_test

import (
	"api/orders/pkg/jwt"
	"testing"
)

func TestJWTCreate(t *testing.T) {
	const phone = "79851174203"
	jwtService := jwt.NewJWT("/2+XnmJGz1j3ehIVI/5P9kl+CghrE3DcS7rnT+qar5w=")
	token, err := jwtService.Create(phone)
	if err != nil {
		t.Fatal(err)
	}
	t_phone, err := jwtService.Parse(token)
	if err != nil {
		t.Fatal(err)
	}
	if t_phone != phone {
		t.Fatalf("Phone %s not equal %s", t_phone, phone)
	}
}
