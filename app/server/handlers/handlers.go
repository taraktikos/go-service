package handlers

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Handler struct {
	Log *zap.Logger
	DB  *sqlx.DB
}

func (h Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("test")
	w.Write([]byte("Go service"))
}
