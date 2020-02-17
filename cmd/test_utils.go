package cmd

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
)

func readFixture(path string) []byte {
	basePath, _ := os.Getwd()
	file := filepath.Join(basePath, path)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return data
}

type simpleMockServer struct {
	server *httptest.Server
}

func NewSimpleMockServer(fixturePath, contentType string, statusCode int) simpleMockServer {
	ts := httptest.NewUnstartedServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				r.Header.Set("Content-Type", contentType)
				b := readFixture(fixturePath)
				w.WriteHeader(statusCode)
				w.Write(b)
			},
		),
	)
	ts.Start()
	return simpleMockServer{server: ts}
}

func (s *simpleMockServer) URL() string {
	return s.server.URL
}

func (s *simpleMockServer) Close() {
	s.server.Close()
}
