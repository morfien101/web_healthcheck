package webserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

//ServerConfig is used to hold the values about the HTTP server.
type ServerConfig struct {
	Cert          string
	Key           string
	ListenAddress string
	Routes        map[string][]byte
}

// Server holds the pointer to the http server.
// It is used to start and stop the HTTP server
type Server struct {
	HTTP   *http.Server
	config *ServerConfig
}

// NewServer will take a ServerConfig, generate a server and return a struct
// containing a pointer to http.Server.
func NewServer(conf *ServerConfig) *Server {
	return &Server{
		HTTP: &http.Server{
			Addr:    conf.ListenAddress,
			Handler: digestRoutes(conf.Routes),
		},
		config: conf,
	}
}

// This is a testing method so that we don't need to fire up the web server.
// It allows us to just test the routes.
func (s *Server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	s.HTTP.Handler.ServeHTTP(w, r)
}

// Start will start the HTTP servers.
// If the configuration has TLS enabled it will be used.
func (s *Server) Start() error {
	if s.config.Cert != "" && s.config.Key != "" {
		fmt.Println("Starting the HTTPS Server")
		return s.HTTP.ListenAndServeTLS(s.config.Cert, s.config.Key)
	}

	fmt.Println("Starting the HTTP Server")
	return s.HTTP.ListenAndServe()
}

// Stop will gracefully try to stop the server. After ten seconds it will
// terminate the server.
func (s *Server) Stop(timeout int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()
	return s.HTTP.Shutdown(ctx)
}

func digestRoutes(routes map[string][]byte) *http.ServeMux {
	r := http.NewServeMux()
	for path, payload := range routes {
		fmt.Printf("Adding route '%s'\n", path)
		r.HandleFunc(
			path,
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("content-type", "application/json")
				w.Write(payload)
			},
		)
	}

	return r
}
