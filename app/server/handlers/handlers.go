package handlers

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Handler struct {
	Log *zap.Logger
	DB  *sqlx.DB
}

func (h Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Tanya"))
}
