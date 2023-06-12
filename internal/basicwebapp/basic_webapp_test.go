package basicwebapp_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/barinek/sonicapp/internal/basicwebapp"
	"github.com/barinek/sonicapp/pkg/testsupport"
	"github.com/barinek/sonicapp/pkg/websupport"
	"github.com/stretchr/testify/assert"
)

func TestNewBasicApp(t *testing.T) {
	withTestServer(func(server *http.Server) {
		resp, _ := http.Get(fmt.Sprintf("http://%s/", server.Addr))
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "Your Application")

		resp, _ = http.Get(fmt.Sprintf("http://%s/terms", server.Addr))
		body, _ = io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "Terms of service")

		resp, _ = http.Get(fmt.Sprintf("http://%s/privacy", server.Addr))
		body, _ = io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "Privacy policy")
	})
}

func withTestServer(f func(server *http.Server)) {

	listener, _ := net.Listen("tcp", "localhost:0")

	basic := basicwebapp.NewBasicApp(nil)
	server := websupport.Create(listener.Addr().String(), basic.LoadHandlers())

	go websupport.Start(server, listener)
	testsupport.WaitForHealthy(server, "health")

	f(server)

	_ = server.Shutdown(context.Background())
}
