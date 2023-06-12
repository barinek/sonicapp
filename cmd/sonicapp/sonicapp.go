package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/barinek/sonicapp/internal/basicwebapp"
	"github.com/barinek/sonicapp/pkg/websupport"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func App(addr string, config *oauth2.Config) *http.Server {
	return websupport.Create(addr, basicwebapp.NewBasicApp(config).LoadHandlers())
}

func newApp(addr string) (*http.Server, net.Listener) {
	if found := os.Getenv("PORT"); found != "" {
		host, _, _ := net.SplitHostPort(addr)
		addr = fmt.Sprintf("%v:%v", host, found)
	}
	log.Printf("Found server address %v", addr)
	listener, _ := net.Listen("tcp", addr)

	clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")

	config := &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
	return App(listener.Addr().String(), config), listener
}

func main() {
	websupport.Start(newApp("0.0.0.0:8888"))
}
