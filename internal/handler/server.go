package handler

import (
	"fmt"
	"net/http"

	"github.com/dawitel/Ashok-reverse-proxy-test/internal/config"
)

type Server struct {
	ProxyHandler func(w http.ResponseWriter, r *http.Request)
	cfg          *config.Config
}

func NewServer(ph func(w http.ResponseWriter, r *http.Request), cfg *config.Config) *Server {
	return &Server{
		ProxyHandler: ph,
		cfg:          cfg,
	}
}

func (s *Server) Start() {
	http.HandleFunc("/", s.ProxyHandler)
	http.HandleFunc("/update-cookie", s.HandleCookieUpdate)
	fmt.Printf("Proxy server is running on :%s\n", s.cfg.Server.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", s.cfg.Server.Port), nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
