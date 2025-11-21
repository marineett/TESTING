package server

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
)

type contextKey string

const (
	ContextUserIDKey   contextKey = "user_id"
	ContextUserRoleKey contextKey = "role"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authz := r.Header.Get("Authorization")
		if !strings.HasPrefix(strings.ToLower(authz), "bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte(ERR_MSG_NO_TOKEN))
			if err != nil {
				http.Error(w, "Error writing data", http.StatusInternalServerError)
				return
			}
			return
		}
		token := strings.TrimSpace(authz[len("Bearer "):])
		claims, err := verifyJWT(token, os.Getenv("JWT_SECRET"))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = w.Write([]byte(ERR_MSG_BAD_TOKEN))
			if err != nil {
				http.Error(w, "Error writing data", http.StatusInternalServerError)
				return
			}
			return
		}
		ctx := r.Context()
		if v, ok := claims["sub"].(float64); ok {
			ctx = context.WithValue(ctx, ContextUserIDKey, int64(v))
		}
		if v, ok := claims["role"].(string); ok {
			ctx = context.WithValue(ctx, ContextUserRoleKey, v)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func verifyJWT(token, secret string) (map[string]interface{}, error) {
	if token == "" || secret == "" {
		return nil, errors.New("empty")
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("fmt")
	}
	head, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	var hdr map[string]interface{}
	if err := json.Unmarshal(head, &hdr); err != nil {
		return nil, err
	}
	if alg, _ := hdr["alg"].(string); alg != "HS256" {
		return nil, errors.New("alg")
	}
	sigInput := parts[0] + "." + parts[1]
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(sigInput))
	expected := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		return nil, errors.New("sig")
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	var claims map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, err
	}
	return claims, nil
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Requested-With")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
