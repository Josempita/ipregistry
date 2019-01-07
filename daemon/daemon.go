package daemon

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Josempita/ipregistry/model"
	"github.com/Josempita/ipregistry/ui"
)

type Config struct {
	ListenSpec string
	UI         ui.Config
}

func Run(cfg *Config) error {
	log.Printf("Starting, HTTP on: %s\n", cfg.ListenSpec)
	m := model.New(nil)
	ui.Start(cfg.UI, m)

	waitForSignal()

	return nil
}

func waitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("Got signal: %v, exiting.", s)
}
