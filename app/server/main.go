package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "embed"

	"github.com/ardanlabs/conf"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/taraktikos/go-service/app/server/handlers"
	"github.com/taraktikos/go-service/business/data/schema"
	"github.com/taraktikos/go-service/foundation/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
)

func main() {

	log := logger.New()
	defer log.Sync()

	if err := run(log); err != nil {
		log.Error("startup", zap.Error(err))
		os.Exit(1)
	}
}

//go:embed static/favicon.ico
var favicon []byte

func run(log *zap.Logger) error {
	cfg := struct {
		Web struct {
			APIHost string `conf:"default:0.0.0.0:8080"`
		}
		DB struct {
			ConnectionString string `conf:"default:postgres://postgres:postgres@localhost:5432/go-service"`
		}
	}{}

	if err := conf.Parse(os.Args[1:], "GO_SERVICE", &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage("GO_SERVICE", &cfg)
			if err != nil {
				return fmt.Errorf("failed to generate config usage: %w", err)
			}
			fmt.Println(usage)
			return nil
		case conf.ErrVersionWanted:
			version, err := conf.VersionString("GO_SERVICE", &cfg)
			if err != nil {
				return fmt.Errorf("failed to generate config version: %w", err)
			}
			fmt.Println(version)
			return nil
		}
		return fmt.Errorf("failed to parse config: %w", err)
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("failed to generate config for output: %w", err)
	}
	log.Info("startup", zap.Any("config", out))

	time.Sleep(20 * time.Second)
	db, err := sqlx.Connect("pgx", cfg.DB.ConnectionString)
	if err != nil {
		log.Error("can't connect to database", zap.Error(err))
		os.Exit(1)
	}

	err = schema.Migrate(context.Background(), db)
	if err != nil {
		log.Error("failed to migrate database", zap.Error(err))
		os.Exit(1)
	}
	log.Info("database is up to date")

	err = schema.Seed(context.Background(), db)
	if err != nil {
		log.Error("failed to seed test data", zap.Error(err))
	}

	hh := handlers.Handler{
		Log: log,
		DB:  db,
	}

	mux := &http.ServeMux{}

	mux.HandleFunc("/", hh.HomePage)
	mux.HandleFunc("/favicon.ico", faviconHander)

	var prod = true
	if prod {
		hostPolicy := func(ctx context.Context, host string) error {
			allowedHost := "new.bankets.com.ua"
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
			log.Info("Starting HTTPS server", zap.String("port", httpsServer.Addr))
			err := httpsServer.ListenAndServeTLS("", "")
			if err != nil {
				log.Error("httpsSrv.ListendAndServeTLS() failed", zap.Error(err))
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
			Addr:         cfg.Web.APIHost,
		}
		log.Info("Starting HTTP server", zap.String("port", cfg.Web.APIHost))
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Error("failed to start server", zap.Error(err))
		}
	} else {
		httpServer := &http.Server{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  120 * time.Second,
			Handler:      mux,
			Addr:         cfg.Web.APIHost,
		}

		log.Info("Starting HTTP server", zap.String("port", cfg.Web.APIHost))
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Error("failed to start server", zap.Error(err))
		}
	}

	return nil
}

func faviconHander(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	w.Write(favicon)
}
