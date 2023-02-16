package mmdbserver

import (
	"encoding/hex"
	"net/http/httptest"
	"os"
	"testing/fstest"
	"time"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	modified := time.Unix(1672502400, 0)
	md5Hash := "0a52730597fb4ffa01fc117d9e71e3a9"
	gzipped := "1f8b08000000000000ff"
	fs := &fstest.MapFS{
		"example.mmdb": {
			Data:    []byte("Example"),
			ModTime: modified,
		},
	}
	var err error
	file, err := fs.Open("example.mmdb")
	assert.NoError(t, err)
	mmdbResponse := NewResponse()
	_, err = mmdbResponse.ReadFile(nil)
	assert.ErrorIs(t, err, ErrDatabaseNotFound)
	_, err = mmdbResponse.ReadFrom(nil)
	assert.ErrorIs(t, err, ErrDatabaseNotFound)
	_, err = mmdbResponse.ReadFile((*os.File)(nil))
	assert.ErrorIs(t, err, ErrDatabaseNotFound)
	_, err = mmdbResponse.ReadFile(file)
	assert.NoError(t, err)
	assert.True(t, mmdbResponse.Valid())
	response := httptest.NewRecorder()
	_, err = mmdbResponse.WriteTo(response)
	assert.NoError(t, err)
	assert.Equal(t, gzipped, hex.EncodeToString(response.Body.Bytes()))
	header := response.Header()
	assert.Equal(t, header.Get("Content-Type"), "application/x-gzip")
	assert.Equal(t, header.Get("Content-Encoding"), "gzip")
	assert.Equal(t, header.Get("Last-Modified"), mmdbResponse.ModTime.Format(time.RFC1123))
	assert.Equal(t, header.Get("X-Database-Md5"), md5Hash)
}
