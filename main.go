package main

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	_ "github.com/kshvakov/clickhouse"
)

func main() {
	cfg, err := _ConfigFromEnv()
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	var logLevel log.Level

	if cfg.Debug {
		logLevel = log.DebugLevel
	} else {
		logLevel = log.InfoLevel
	}

	log.SetLevel(logLevel)

	log.Debugf("max workers size is %d", cfg.MaxWorkers)

	args := append([]string{"clickhouse-client"}, os.Args[1:]...)

	query := strings.Join(args, " ")

	log.Debugf("searches shodan using \"%s\"", query)

	cd, err := _NewClickDown(cfg)
	if err != nil {
		log.Fatalf("failed to init clickdown: %v", err)
	}
	defer func() {
		if err := cd.Shutdown(); err != nil {
			log.Fatalf("failed to shutdown: %v", err)
		}
	}()

	if err := cd.Run(query); err != nil {
		log.Fatalf("failed to run clickdown: %v", err)
	}
}
