package mmdbserver

import (
	"bufio"
	"net/http"
	"os"
	"path/filepath"
)

func mustFixtureRequest(name string) (request *http.Request) {
	fp, err := os.Open(filepath.Join("fixtures", name, "request.http"))
	if err != nil {
		panic(err)
	}
	request, err = http.ReadRequest(bufio.NewReader(fp))
	if err != nil {
		panic(err)
	}
	return
}

func mustFixtureResponse(name string, request *http.Request) (response *http.Response) {
	fp, err := os.Open(filepath.Join("fixtures", name, "response.http"))
	if err != nil {
		panic(err)
	}
	response, err = http.ReadResponse(bufio.NewReader(fp), request)
	if err != nil {
		panic(err)
	}
	return
}
