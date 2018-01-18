package discovery

import (
	"fmt"

	"github.com/ernoaapa/eliot/pkg/version"
	"github.com/grandcat/zeroconf"
	log "github.com/sirupsen/logrus"
)

// Server is zeroconf discovery server
type Server struct {
	Name     string
	Domain   string
	Port     int
	server   *zeroconf.Server
	shutdown chan bool
}

// NewServer creates new discovery server
func NewServer(name string, port int) *Server {
	return &Server{
		Name:   name,
		Domain: "local.",
		Port:   port,
	}
}

// Serve starts the discovery server
func (s *Server) Serve() {
	log.Infof("Start discovery server...")
	log.Debugf("Exposing %s in port %d", s.Name, s.Port)
	server, err := zeroconf.Register(s.Name, ZeroConfServiceName, s.Domain, s.Port, []string{
		fmt.Sprintf("v=%s", version.VERSION),
	}, nil)
	if err != nil {
		log.Fatalf("Failed to create zeroconf server: %s", err)
	}

	s.server = server

	select {
	case <-s.shutdown:
		s.server.Shutdown()
	}
}

// Stop server to be discoverable
func (s *Server) Stop() {
	log.Infof("Stop discovery server...")
	s.shutdown <- true
}
