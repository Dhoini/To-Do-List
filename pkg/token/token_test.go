package token

import "testing"

const email = "test@gmail.com"

func TestJWTSecret_GenerateToken(t *testing.T) {
	jwtService := NewJWT("RxbxgRcFCFes0enila83XSdWzejBmKuw4cHiPuMgiU8")

	token, err := jwtService.GenerateToken(JwtDate{
		Email: email,
	})
	if err != nil {
		t.Fatalf("Error generating token: %v", err)
		return
	}

	isValid, data := jwtService.ParseToken(token)
	if !isValid || data == nil {
		t.Fatalf("Error validating token: %v", data)
		return
	}

	if data.Email != email {
		t.Fatalf("Error validating token: %v", data.Email)
		return
	}

}
