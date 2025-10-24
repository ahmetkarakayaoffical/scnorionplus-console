package authserver

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"

	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/controllers/authserver/handlers"
	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/controllers/router"
	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/controllers/sessions"
	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/models"
	"github.com/ahmetkarakayaoffical/scnorionplus-console/utils"
	"github.com/labstack/echo/v4"
)

type AuthServer struct {
	Router         *echo.Echo
	Handler        *handlers.Handler
	Server         *http.Server
	SessionManager *sessions.SessionManager
	CACert         *x509.Certificate
}

func New(m *models.Model, s *sessions.SessionManager, caCert, server, consolePort, authPort, reverseProxyAuthPort string) *AuthServer {
	var err error
	a := AuthServer{}

	// Get max upload size setting
	maxUploadSize, err := m.GetMaxUploadSize()
	if err != nil {
		maxUploadSize = "512M"
		log.Println("[ERROR]: could not get max upload size from database")
	}

	// Router
	a.Router = router.New(s, server, authPort, maxUploadSize)

	// Session Manager
	a.SessionManager = s

	a.CACert, err = utils.ReadPEMCertificate(caCert)
	if err != nil {
		log.Fatal(err)
	}

	// Create Handlers and register its router
	a.Handler = handlers.NewHandler(m, s, a.CACert, server, consolePort, reverseProxyAuthPort)
	a.Handler.Register(a.Router)

	return &a
}

func (a *AuthServer) Serve(address, certFile, certKey string) error {
	cp := x509.NewCertPool()
	cp.AddCert(a.CACert)

	a.Server = &http.Server{
		Addr:    address,
		Handler: a.Router,
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequestClientCert,
			ClientCAs:  cp,
		},
	}
	return a.Server.ListenAndServeTLS(certFile, certKey)
}

func (a *AuthServer) Close() error {
	return a.Server.Close()
}
