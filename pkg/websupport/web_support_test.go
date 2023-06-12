package websupport_test

import (
	"net"
	"testing"

	"github.com/barinek/sonicapp/pkg/healthsupport"
	"github.com/barinek/sonicapp/pkg/testsupport"
	"github.com/barinek/sonicapp/pkg/websupport"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	listener, _ := net.Listen("tcp", "localhost:0")
	server := websupport.Create(listener.Addr().String(), func(router *mux.Router) {
		router.HandleFunc("/health", healthsupport.HandlerFunction).Methods("GET")
	})
	go websupport.Start(server, listener)
	assert.True(t, testsupport.WaitForHealthy(server, "health"))
	websupport.Stop(server)
}
