package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pdkonovalov/auth-server/pkg/config"
	"github.com/pdkonovalov/auth-server/pkg/jwt"
	"github.com/pdkonovalov/auth-server/pkg/storage/postgres"
)

func TestNewJwtEndpoint(t *testing.T) {
	config, err := config.ReadConfig(os.Getenv)
	if err != nil {
		t.Fatalf("error load config: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go run(ctx, os.Getenv)
	time.Sleep(2 * time.Second)

	u := url.URL{
		Scheme:   "http",
		Host:     net.JoinHostPort(config.Host, config.Port),
		Path:     "/api/v1/jwt/new",
		RawQuery: "guid=" + uuid.New().String(),
	}
	httpResponse, err := http.Get(u.String())
	if err != nil || httpResponse.StatusCode != http.StatusOK {
		t.Fatalf("cant get new access and refresh tokens: %s", err)
	}
	response := &struct {
		AccessToken  string
		RefreshToken string
	}{}
	d := json.NewDecoder(httpResponse.Body)
	err = d.Decode(&response)
	if err != nil {
		t.Fatalf("invalid response format: %s", err)
	}

	jwt, err := jwt.Init(config)
	if err != nil {
		t.Fatalf("error init jwt: %s", err)
	}
	_, jtiAccess, valid := jwt.ValidateAccessToken(response.AccessToken)
	if !valid {
		t.Fatal("invalid access token")
	}
	jtiRefresh, valid := jwt.ValidateRefreshToken(response.RefreshToken)
	if !valid {
		t.Fatal("invalid refresh token")
	}
	if jtiAccess != jtiRefresh {
		t.Fatal("access and refresh token have different jti")
	}

	storage, err := postgres.Init(config)
	if err != nil {
		t.Fatalf("error init storage: %s", err)
	}

	defer storage.Shutdown()

	err = storage.DeleteJti(jtiAccess)
	if err != nil {
		t.Fatalf("error clear test data from storage: %s", err)
	}
}

func TestRefreshJwtEndpoint(t *testing.T) {
	config, err := config.ReadConfig(os.Getenv)
	if err != nil {
		t.Fatalf("error load config: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go run(ctx, os.Getenv)
	time.Sleep(2 * time.Second)

	u := url.URL{
		Scheme:   "http",
		Host:     net.JoinHostPort(config.Host, config.Port),
		Path:     "/api/v1/jwt/new",
		RawQuery: "guid=" + uuid.New().String(),
	}
	httpResponse, err := http.Get(u.String())
	if err != nil || httpResponse.StatusCode != http.StatusOK {
		t.Fatalf("cant get new access and refresh tokens: %s", err)
	}
	response := &struct {
		AccessToken  string
		RefreshToken string
	}{}
	d := json.NewDecoder(httpResponse.Body)
	err = d.Decode(&response)
	if err != nil {
		t.Fatalf("invalid response format: %s", err)
	}

	jwt, err := jwt.Init(config)
	if err != nil {
		t.Fatalf("error init jwt: %s", err)
	}
	_, jtiAccess, valid := jwt.ValidateAccessToken(response.AccessToken)
	if !valid {
		t.Fatal("invalid access token")
	}
	jtiRefresh, valid := jwt.ValidateRefreshToken(response.RefreshToken)
	if !valid {
		t.Fatal("invalid refresh token")
	}
	if jtiAccess != jtiRefresh {
		t.Fatal("access and refresh token have different jti")
	}

	oldAccessToken := response.AccessToken
	oldRefreshToken := response.RefreshToken
	oldJti := jtiAccess

	request := &struct {
		AccessToken  string
		RefreshToken string
	}{oldAccessToken, oldRefreshToken}
	requestBody, _ := json.Marshal(request)
	u = url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(config.Host, config.Port),
		Path:   "/api/v1/jwt/refresh",
	}
	httpResponse, err = http.Post(u.String(), "application/json", bytes.NewBuffer(requestBody))
	if err != nil || httpResponse.StatusCode != http.StatusOK {
		t.Fatalf("cant refresh tokens: %s", err)
	}

	d = json.NewDecoder(httpResponse.Body)
	err = d.Decode(&response)
	if err != nil {
		t.Fatalf("invalid response format: %s", err)
	}
	_, jtiAccess, valid = jwt.ValidateAccessToken(response.AccessToken)
	if !valid {
		t.Fatal("invalid access token")
	}
	jtiRefresh, valid = jwt.ValidateRefreshToken(response.RefreshToken)
	if !valid {
		t.Fatal("invalid refresh token")
	}
	if jtiAccess != jtiRefresh {
		t.Fatal("access and refresh token have different jti")
	}

	newAccessToken := response.AccessToken
	newRefreshToken := response.RefreshToken
	newJti := jtiAccess

	if newAccessToken == oldAccessToken {
		t.Fatal("access token not refreshed")
	}
	if newRefreshToken == oldRefreshToken {
		t.Fatal("refresh token not refreshed")
	}
	if newJti == oldJti {
		t.Fatal("new tokens have same jti")
	}

	request = &struct {
		AccessToken  string
		RefreshToken string
	}{oldAccessToken, oldRefreshToken}
	requestBody, _ = json.Marshal(request)
	httpResponse, _ = http.Post(u.String(), "application/json", bytes.NewBuffer(requestBody))
	if httpResponse.StatusCode == http.StatusOK {
		t.Fatal("refresh operation for old refresh token with success")
	}

	storage, err := postgres.Init(config)
	if err != nil {
		t.Fatalf("error init storage: %s", err)
	}

	defer storage.Shutdown()

	err = storage.DeleteJti(oldJti)
	if err != nil {
		t.Fatalf("error clear test data from storage: %s", err)
	}
	err = storage.DeleteJti(newJti)
	if err != nil {
		t.Fatalf("error clear test data from storage: %s", err)
	}
}
