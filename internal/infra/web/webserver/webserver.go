package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

// AddHandler adds a new handler to the WebServer.
func (s *WebServer) AddHandler(path, method string, handler http.HandlerFunc) {
	if s.Handlers[path] == nil {
		s.Handlers[path] = make(map[string]http.HandlerFunc)
	}
	s.Handlers[path][method] = handler
}

// Start loops through the handlers and adds them to the router.
// register middeleware logger
// start the server
func (s *WebServer) Start() {
	for path, methods := range s.Handlers {
		for method, handler := range methods {
			s.Router.MethodFunc(method, path, handler)
		}
	}
	http.ListenAndServe(s.WebServerPort, s.Router)
}
