package client

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"alphatr.com/acme-lego/src/config"
	"github.com/go-acme/lego/v3/lego"
)

type HTTPProviderServer struct {
	iface    string
	port     string
	done     chan bool
	listener net.Listener
}

func init() {
	ProviderMap["http-port"] = ApplyHTTPPortProvider
}

func ApplyHTTPPortProvider(domain string, cli *lego.Client, conf *config.DomainConfig) error {
	host, port, err := net.SplitHostPort(":8057")
	if err != nil {
		return err
	}

	return cli.Challenge.SetHTTP01Provider(NewHTTPProviderServer(host, port))
}

func HTTP01ChallengePath(token string) string {
	return "/.well-known/acme-challenge/" + token
}

func NewHTTPProviderServer(iface, port string) *HTTPProviderServer {
	return &HTTPProviderServer{iface: iface, port: port}
}

func (s *HTTPProviderServer) Present(domain, token, keyAuth string) error {
	if s.port == "" {
		s.port = "80"
	}

	var err error
	s.listener, err = net.Listen("tcp", net.JoinHostPort(s.iface, s.port))
	if err != nil {
		return fmt.Errorf("Could not start HTTP server for challenge -> %v", err)
	}

	s.done = make(chan bool)
	go s.serve(domain, token, keyAuth)
	return nil
}

func (s *HTTPProviderServer) CleanUp(domain, token, keyAuth string) error {
	if s.listener == nil {
		return nil
	}
	s.listener.Close()
	<-s.done
	return nil
}

func (s *HTTPProviderServer) serve(domain, token, keyAuth string) {
	path := HTTP01ChallengePath(token)

	mux := http.NewServeMux()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.Host, domain) && r.Method == "GET" {
			w.Header().Add("Content-Type", "text/plain")
			w.Write([]byte(keyAuth))
		} else {
			w.Write([]byte("TEST"))
		}
	})

	httpServer := &http.Server{
		Handler: mux,
	}

	httpServer.SetKeepAlivesEnabled(false)
	httpServer.Serve(s.listener)
	s.done <- true
}
