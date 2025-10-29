package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"data_base_project/service_logic"
	"data_base_project/types"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func SetupAuthorizeRouterV2(authService service_logic.IAuthService) *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("login", AuthorizeHandlerV2(authService)).Methods("POST")
	return router
}

func AuthorizeHandlerV2(authService service_logic.IAuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		var authData types.ServerAuthData
		if err := json.Unmarshal(body, &authData); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}
		verdict, err := authService.Authorize(*types.MapperAuthServerToService(&authData))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			http.Error(w, "JWT secret is not configured", http.StatusInternalServerError)
			return
		}
		token, err := createJWT(verdict.UserID, verdict.UserType.String(), 24*time.Hour, secret)
		if err != nil {
			http.Error(w, "Failed to issue token", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		resp := types.ServerAuthResponseV2{Token: token, Role: verdict.UserType.String(), UserID: verdict.UserID}
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func SetupAuthorizeRouterV1(authService service_logic.IAuthService, logger *log.Logger) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc(AUTH_AUTHORIZE, AuthorizeHandler(authService, logger))

	return router
}

func AuthorizeHandler(authService service_logic.IAuthService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "POST" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Printf("Error reading request body: %v", err)
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		var authData types.ServerAuthData
		if err := json.Unmarshal(body, &authData); err != nil {
			logger.Printf("Error unmarshaling request body: %v", err)
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}
		verdict, err := authService.Authorize(*types.MapperAuthServerToService(&authData))
		if err != nil {
			logger.Printf("Error authorizing: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		serverVerdict := types.MapperVerdictServiceToServer(&verdict)
		json.NewEncoder(w).Encode(serverVerdict)
		logger.Printf("Authorized: %v", verdict)
	}
}

// createJWT issues an HS256 JWT compatible with JWTAuthMiddleware
func createJWT(sub int64, role string, ttl time.Duration, secret string) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	now := time.Now().Unix()
	payload := map[string]interface{}{
		"sub":  sub,
		"role": role,
		"iat":  now,
		"exp":  now + int64(ttl.Seconds()),
	}
	hb, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	pb, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	hEnc := base64.RawURLEncoding.EncodeToString(hb)
	pEnc := base64.RawURLEncoding.EncodeToString(pb)
	sigInput := hEnc + "." + pEnc
	sig := hmacSign(sigInput, secret)
	token := sigInput + "." + sig
	return token, nil
}

func hmacSign(input, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(input))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
