package main

import (
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/taraktikos/go-service/foundation/logger"
	"go.uber.org/zap"
	"os"
)

func main() {
	fmt.Println("Server")
	// Perform the startup and shutdown sequence.
	log := logger.New("SALES-API")
	//defer log.Sync()

	if err := run(log); err != nil {
		//log.Errorw("startup", "ERROR", err)
		os.Exit(1)
	}
}

func run(log *zap.Logger) error {
	cfg := struct {
		conf.Version
	}{}

	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				//return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		case conf.ErrVersionWanted:
			version, err := conf.VersionString("SALES", &cfg)
			if err != nil {
				//return errors.Wrap(err, "generating config version")
			}
			fmt.Println(version)
			return nil
		}
		//return errors.Wrap(err, "parsing config")
	}

	return nil
}
