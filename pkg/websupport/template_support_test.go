package websupport_test

import (
	"bytes"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/barinek/sonicapp/pkg/websupport"
	"github.com/barinek/sonicapp/pkg/websupport/test"
	"github.com/stretchr/testify/assert"
)

func TestModelAndView(t *testing.T) {
	w := &httptest.ResponseRecorder{Body: new(bytes.Buffer)}
	data := websupport.Model{Map: map[string]any{"name": "aName", "float": 3.1415, "percentage": .8843}}
	_ = websupport.ModelAndView(w, &websupport_test.Resources, data, "test")
	body, _ := io.ReadAll(w.Body)

	assert.Equal(t, `
    <!DOCTYPE html>
    <html lang="en">
    <body>aName 3 88%
    </body>
    </html>`, string(body))
}
