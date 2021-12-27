package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/braintree/manners"
	"github.com/caarlos0/env"
	"github.com/kyxap1/pwgen/handlers"
	"github.com/kyxap1/pwgen/health"
)

const version = "1.0.0"

type config struct {
	HttpHost string `env:"HTTP_HOST" envDefault:"0.0.0.0"`
	HttpPort string `env:"HTTP_PORT" envDefault:"8080"`

	PasswordLength int  `env:"PASSWORD_LENGTH" envDefault:"16"`
	NumDigits      int  `env:"NUM_DIGITS" envDefault:"1"`
	NumSymbols     int  `env:"NUM_SYMBOLS" envDefault:"1"`
	NoUpper        bool `env:"NO_UPPER" envDefault:"false"`
	AllowRepeat    bool `env:"ALLOW_REPEAT" envDefault:"false"`
}

func main() {
	log.Println("Starting application version:", version)

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	httpAddr := fmt.Sprintf("%s:%s", cfg.HttpHost, cfg.HttpPort)
	log.Printf("HTTP service listening on %s", httpAddr)

	mux := http.NewServeMux()
	mux.Handle("/", handlers.PwgenHandler(cfg.PasswordLength, cfg.NumDigits,
		cfg.NumSymbols, cfg.NoUpper, cfg.AllowRepeat))
	mux.HandleFunc("/healthz", health.HealthzHandler)
	mux.HandleFunc("/healthz/status", health.HealthzStatusHandler)
	mux.Handle("/version", handlers.VersionHandler(version))

	httpServer := manners.NewServer()
	httpServer.Addr = httpAddr
	httpServer.Handler = handlers.LoggingHandler(mux)

	errChan := make(chan error, 10)

	go func() {
		errChan <- httpServer.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))
			httpServer.BlockingClose()
			os.Exit(0)
		}
	}
}
