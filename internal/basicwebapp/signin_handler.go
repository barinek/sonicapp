package basicwebapp

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"net/http"

	"github.com/barinek/sonicapp/pkg/websupport"
	"github.com/coreos/go-oidc/v3/oidc"
)

func (a BasicApp) signIn(writer http.ResponseWriter, request *http.Request) {
	model := websupport.Model{Map: map[string]any{}}
	_ = websupport.ModelAndView(writer, &Resources, model, "sign-in")
}

func (a BasicApp) authenticate(writer http.ResponseWriter, request *http.Request) {
	session, _ := a.store.Get(request, a.cookieName)

	state := a.randomString()
	session.Values["state"] = state
	err := session.Save(request, writer)
	if err != nil {
		http.Error(writer, "oops", http.StatusInternalServerError)
		return
	}
	authCodeURL := a.config.AuthCodeURL(state)

	http.Redirect(writer, request, authCodeURL, http.StatusFound)
}

func (a BasicApp) randomString() string {
	bytes := make([]byte, 16)
	_, _ = io.ReadFull(rand.Reader, bytes)
	state := base64.RawURLEncoding.EncodeToString(bytes)
	return state
}

func (a BasicApp) callback(writer http.ResponseWriter, request *http.Request) {
	state := request.URL.Query().Get("state")

	session, err := a.store.Get(request, a.cookieName)
	if err != nil {
		http.Error(writer, "oops", http.StatusInternalServerError)
	}

	session.Options.MaxAge = -1
	_ = session.Save(request, writer)

	if session.Values["state"] != state {
		http.Error(writer, "oops", http.StatusInternalServerError)
		return
	}

	code := request.URL.Query().Get("code")
	oauth2Token, err := a.config.Exchange(context.Background(), code)
	rawToken, _ := oauth2Token.Extra("id_token").(string)
	provider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	verifier := provider.Verifier(&oidc.Config{ClientID: a.config.ClientID})
	idToken, err := verifier.Verify(context.Background(), rawToken)

	claims := struct {
		Email      string `json:"email,omitempty"`
		GivenName  string `json:"given_name,omitempty"`
		FamilyName string `json:"family_name,omitempty"`
	}{}
	err = idToken.Claims(&claims)

	if err != nil {
		log.Println(err)

		http.Error(writer, "oops", http.StatusInternalServerError)
	}

	session.Values["principal"] = claims.Email
	session.Options.MaxAge = a.maxAge
	_ = session.Save(request, writer)

	http.Redirect(writer, request, "/", http.StatusFound)
}

func (a BasicApp) logout(writer http.ResponseWriter, request *http.Request) {
	session, _ := a.store.Get(request, a.cookieName)
	session.Values = nil
	session.Options.MaxAge = -1
	err := session.Save(request, writer)
	if err != nil {
		return
	}
	http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
}
