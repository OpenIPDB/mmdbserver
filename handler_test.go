package mmdbserver

import (
	"bytes"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMMDBUpdateHandler(t *testing.T) {
	modified := time.Unix(1672502400, 0)

	handler := &MMDBUpdateHandler{
		HomePage: "https://example.com",
		MMDBHandler: MMDBHandlerFunc(func(response *Response, request *Request) (err error) {
			response.ModTime = modified
			if strings.HasPrefix(request.EditionID, "GeoLite2") {
				var payload bytes.Buffer
				payload.WriteString("example")
				_, _ = response.ReadFrom(&payload)
			} else if request.EditionID == "Unexpected" {
				err = io.ErrUnexpectedEOF
			}
			return
		}),
	}

	fixtures := []string{
		"general",
		"invalid-edition-id",
		"invalid-account-id",
		"invalid-license-key",
		"invalid-hash",
		"method-not-allowed",
		"internal-server-error",
		"database-not-found",
		"database-latest",
	}

	for _, name := range fixtures {
		request := mustFixtureRequest(name)
		actualResponse := httptest.NewRecorder()
		handler.ServeHTTP(actualResponse, request)
		expectedResponse := mustFixtureResponse(name, request)
		assert.Equal(t, expectedResponse.StatusCode, actualResponse.Code, name)
		assert.Equal(t, expectedResponse.Header, actualResponse.Header(), name)
	}
}
