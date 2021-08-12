package web

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
)

func StartHttpServer(mux http.Handler, addr string, log *zap.Logger) {
	httpServer := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
		Addr:         addr,
	}

	log.Info("starting http server", zap.String("addr", addr))
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Error("failed to start server", zap.Error(err))
	}
}

func StartHttpsServer(mux http.Handler, addr string, log *zap.Logger, domain string) {
	hostPolicy := func(ctx context.Context, host string) error {
		allowedHost := domain
		if host == allowedHost {
			return nil
		}
		return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
	}

	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
		Cache:      autocert.DirCache("."),
	}
	httpsServer := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
		Addr:         ":443",
		TLSConfig:    &tls.Config{GetCertificate: m.GetCertificate},
	}
	go func() {
		log.Info("starting https server", zap.String("addr", httpsServer.Addr))
		err := httpsServer.ListenAndServeTLS("", "")
		if err != nil {
			log.Error("failed to start https server", zap.Error(err))
		}
	}()

	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		newURI := "https://" + r.Host + r.URL.String()
		http.Redirect(w, r, newURI, http.StatusFound)
	}
	muxRedirect := &http.ServeMux{}
	muxRedirect.HandleFunc("/", handleRedirect)

	httpServer := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      m.HTTPHandler(muxRedirect),
		Addr:         addr,
	}
	log.Info("starting http server", zap.String("addr", addr))
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Error("failed to start server", zap.Error(err))
	}
}
