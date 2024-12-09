package v1

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/pdkonovalov/auth-server/pkg/email"
	"github.com/pdkonovalov/auth-server/pkg/jwt"
	"github.com/pdkonovalov/auth-server/pkg/storage"
)

func HandleNewJwt(storage storage.Storage, jwt *jwt.JwtGenerator) http.HandlerFunc {
	type response struct {
		AccessToken  string `json:"AccessToken"`
		RefreshToken string `json:"RefreshToken"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		guid := r.URL.Query().Get("guid")
		_, err := uuid.Parse(guid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		jti, err := storage.WriteNewJti(guid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ip := r.RemoteAddr
		accessToken, err := jwt.GenerateAccessToken(ip, jti)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		refreshToken, err := jwt.GenerateRefreshToken(jti)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp, err := json.Marshal(&response{accessToken, refreshToken})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func HandleRefreshJwt(storage storage.Storage, email *email.Email, jwt *jwt.JwtGenerator) http.HandlerFunc {
	type request struct {
		AccessToken  string `json:"AccessToken"`
		RefreshToken string `json:"RefreshToken"`
	}
	type response struct {
		AccessToken  string `json:"AccessToken"`
		RefreshToken string `json:"RefreshToken"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ip, jtiAccess, valid := jwt.ValidateAccessToken(req.AccessToken)
		if !valid {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		jtiRefresh, valid := jwt.ValidateRefreshToken(req.RefreshToken)
		if !valid {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if jtiAccess != jtiRefresh {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		jti := jtiAccess
		guid, exist, err := storage.FindJti(jti)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !exist {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if ip != r.RemoteAddr {
			email.SendAllert(guid)
		}
		err = storage.DeleteJti(jti)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		newJti, err := storage.WriteNewJti(guid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ip = r.RemoteAddr
		accessToken, err := jwt.GenerateAccessToken(ip, newJti)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		refreshToken, err := jwt.GenerateRefreshToken(newJti)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp, err := json.Marshal(&response{accessToken, refreshToken})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}
