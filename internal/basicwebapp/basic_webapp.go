package basicwebapp

import (
	"io/fs"
	"net/http"

	"github.com/barinek/sonicapp/pkg/healthsupport"
	"github.com/barinek/sonicapp/pkg/metricssupport"
	"github.com/barinek/sonicapp/pkg/websupport"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"golang.org/x/oauth2"
)

type BasicApp struct {
	config     *oauth2.Config
	key        string
	store      sessions.Store
	cookieName string
	maxAge     int
}

func NewBasicApp(config *oauth2.Config) BasicApp {
	return BasicApp{
		config:     config,
		store:      sessions.NewCookieStore([]byte("supersecret")),
		cookieName: "openai",
		maxAge:     60 * 60 * 2,
	}
}

func (a BasicApp) LoadHandlers() func(x *mux.Router) {
	return func(router *mux.Router) {
		router.HandleFunc("/", a.dashboard).Methods("GET")

		router.HandleFunc("/sign-in", a.signIn).Methods("GET")
		router.HandleFunc("/authenticate", a.authenticate).Methods("GET")
		router.HandleFunc("/oauth/callback", a.callback).Methods("GET")
		router.HandleFunc("/logout", a.logout).Methods("GET")

		router.HandleFunc("/health", healthsupport.HandlerFunction)
		router.HandleFunc("/metrics", metricssupport.HandlerFunction)
		router.HandleFunc("/terms", a.terms).Methods("GET")
		router.HandleFunc("/privacy", a.privacy).Methods("GET")
		router.Use(metricssupport.Middleware)

		static, _ := fs.Sub(Resources, "resources/static")
		fileServer := http.FileServer(http.FS(static))
		router.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))
	}
}

func (a BasicApp) dashboard(writer http.ResponseWriter, request *http.Request) {
	session, _ := a.store.Get(request, a.cookieName)

	if session.Values["principal"] == nil {
		http.Redirect(writer, request, "sign-in", http.StatusFound)
		return
	}

	model := websupport.Model{Map: map[string]any{
		"principal": a.principal(request),
	}}
	_ = websupport.ModelAndView(writer, &Resources, model, "index")
}

func (a BasicApp) terms(writer http.ResponseWriter, request *http.Request) {
	model := websupport.Model{Map: map[string]any{
		"principal": a.principal(request),
	}}
	_ = websupport.ModelAndView(writer, &Resources, model, "terms")
}

func (a BasicApp) privacy(writer http.ResponseWriter, request *http.Request) {
	model := websupport.Model{Map: map[string]any{
		"principal": a.principal(request),
	}}
	_ = websupport.ModelAndView(writer, &Resources, model, "privacy")
}

func (a BasicApp) principal(req *http.Request) interface{} {
	session, _ := a.store.Get(req, a.cookieName)
	if session.IsNew {
		return nil
	}
	return session.Values["principal"]
}
