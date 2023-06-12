package metricssupport_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/barinek/sonicapp/pkg/metricssupport"
	"github.com/barinek/sonicapp/pkg/testsupport"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/metrics", metricssupport.HandlerFunction).Methods("GET")
	server := testsupport.Server(r)
	testsupport.WaitForHealthy(server, "metrics")

	resp, _ := http.Get(fmt.Sprintf("http://%s/metrics", server.Addr))
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "{\"requests\":\"0\"}", string(body))

	_ = server.Shutdown(context.Background())
}

func TestMiddleware(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/metrics", metricssupport.HandlerFunction).Methods("GET")
	server := testsupport.Server(r)
	server.Handler.(*mux.Router).Use(metricssupport.Middleware)
	testsupport.WaitForHealthy(server, "metrics")

	resp, _ := http.Get(fmt.Sprintf("http://%s/metrics", server.Addr))
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "{\"requests\":\"2\"}", string(body))

	_ = server.Shutdown(context.Background())
}
