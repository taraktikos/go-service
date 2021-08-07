package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

func New() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("can't create logger instance")
		os.Exit(1)
	}
	return logger
}
