package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/namsral/flag"
	"github.com/pkg/errors"

	"myst/pkgs/logger"
	"myst/server/apis"
	"myst/server/core"
	"myst/server/storage"
)

var log = logger.New("app", logger.Red)

type Config struct {
	Port          int
	DatabasePath  string
	JWTSigningKey string
}

func parseFlags() Config {
	cfg := Config{}

	flag.IntVar(&cfg.Port, "port", 8080, "Port the server should listen on")

	flag.StringVar(&cfg.DatabasePath, "db", "data/myst.sqlite", "SQLite database file")

	flag.StringVar(&cfg.JWTSigningKey, "jwt-signing-key", "", "The key used for signing JWT tokens")

	flag.Parse()

	return cfg
}

func main() {
	cleanup, err := run()
	if err != nil {
		log.Error(err)
		os.Exit(1)
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	err = cleanup()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func run() (cleanup func() error, err error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, errors.WithMessage(err, "could not load .env")
	}

	cfg := parseFlags()

	jwtSigningKey, err := base64.StdEncoding.DecodeString(cfg.JWTSigningKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode jwt signing key")
	}

	logger.EnableDebug = true

	repo, err := storage.Open(cfg.DatabasePath)
	if err != nil {
		return nil, errors.WithMessage(err, "could not open SQLite database")
	}

	app := core.New(
		core.WithKeystoreRepository(repo),
		core.WithUserRepository(repo),
		core.WithInvitationRepository(repo),
	)

	server := apis.NewServer(app, jwtSigningKey)

	err = server.Start(fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		_ = repo.Close()
		return nil, errors.WithMessage(err, "could not start server")
	}

	return func() error {
		stopErr := server.Stop()
		closeErr := repo.Close()
		if stopErr != nil {
			return errors.WithMessage(stopErr, "could not stop server")
		}
		return closeErr
	}, nil
}
