package challenge

import (
	"net"
	"net/http"
	"strings"

	"github.com/go-acme/lego/v3/challenge"

	"github.com/alphatr/acme-lego/common/config"
	"github.com/alphatr/acme-lego/common/errors"
)

// HTTPProviderServer HTTP 端口转发服务器
type HTTPProviderServer struct {
	iface    string
	port     string
	done     chan bool
	listener net.Listener
}

func init() {
	ProviderMap["http-port"] = &HTTPPortProvider{isHTTPS: false}
	ProviderMap["https-port"] = &HTTPPortProvider{isHTTPS: true}
}

// HTTPPortProvider HTTPPortProvider
type HTTPPortProvider struct {
	isHTTPS bool
}

// Type 返回注册的类型
func (ins *HTTPPortProvider) Type() ProviderType {
	if ins.isHTTPS {
		return ProviderTLS
	}

	return ProviderHTTP
}

// Provider Provider 实体
func (ins *HTTPPortProvider) Provider(domain string, conf *config.DomainConf) (challenge.Provider, *errors.Error) {
	host, port, err := net.SplitHostPort(conf.Options["server"])
	if err != nil {
		return nil, errors.NewError(errors.CommonParseHostPortErrno, err)
	}

	provider := NewHTTPProviderServer(host, port)
	return provider, nil
}

// NewHTTPProviderServer 创建端口转发服务器
func NewHTTPProviderServer(iface, port string) *HTTPProviderServer {
	return &HTTPProviderServer{iface: iface, port: port}
}

// Present 启动服务器
func (s *HTTPProviderServer) Present(domain, token, keyAuth string) error {
	if s.port == "" {
		s.port = "80"
	}

	var err error
	s.listener, err = net.Listen("tcp", net.JoinHostPort(s.iface, s.port))
	if err != nil {
		return errors.NewError(errors.ModelChalServerStartErrno, err)
	}

	s.done = make(chan bool)
	go s.serve(domain, token, keyAuth)
	return nil
}

// CleanUp 关闭服务器
func (s *HTTPProviderServer) CleanUp(domain, token, keyAuth string) error {
	if s.listener == nil {
		return nil
	}

	s.listener.Close()
	<-s.done
	return nil
}

func (s *HTTPProviderServer) serve(domain, token, keyAuth string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/acme-challenge/"+token, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.Host, domain) && r.Method == "GET" {
			w.Header().Add("Content-Type", "text/plain")
			w.Write([]byte(keyAuth))
		}
	})

	httpServer := &http.Server{
		Handler: mux,
	}

	httpServer.SetKeepAlivesEnabled(false)
	httpServer.Serve(s.listener)
	s.done <- true
}
