package main

import (
	"context"
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
	"github.com/taraktikos/go-service/foundation/web"
	"go.uber.org/zap"
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
			SSL struct {
				Enabled bool   `conf:"default:false"`
				FQDN    string `conf:"default:go-service.com"`
			}
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

	if cfg.Web.SSL.Enabled {
		web.StartHttpsServer(mux, cfg.Web.APIHost, log, cfg.Web.SSL.FQDN)
	} else {
		web.StartHttpServer(mux, cfg.Web.APIHost, log)
	}

	return nil
}

func faviconHander(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	w.Write(favicon)
}
