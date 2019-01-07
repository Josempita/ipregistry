// Copyright 2015 Daniel Theophanes.
// Use of this source code is governed by a zlib-style
// license that can be found in the LICENSE file.

// Simple service that only works by printing a log message every few seconds.
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Josempita/ipregistry/daemon"

	"github.com/kardianos/service"
)

var logger service.Logger

// Program structures.
//  Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager.")
	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() error {
	cfg := processFlags()

	setupHttpAssets(cfg)

	if err := daemon.Run(cfg); err != nil {
		log.Printf("Error in main(): %v", err)
	}
	return nil
}

var assetsPath string

func processFlags() *daemon.Config {
	cfg := &daemon.Config{}

	flag.StringVar(&cfg.ListenSpec, "listen", "0.0.0.0:4000", "HTTP listen spec")
	flag.StringVar(&assetsPath, "assets-path", "assets", "Path to assets dir")

	flag.Parse()
	return cfg
}

func setupHttpAssets(cfg *daemon.Config) {
	log.Printf("Assets served from %q.", assetsPath)
	cfg.UI.Assets = http.Dir(assetsPath)
	cfg.UI.Host = "0.0.0.0"
	cfg.UI.Port = "4000"
}
func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Info("I'm Stopping!")
	close(p.exit)
	return nil
}

// Service setup.
//   Define service config.
//   Create the service.
//   Setup the logger.
//   Handle service controls (optional).
//   Run the service.
func main() {

	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	svcConfig := &service.Config{
		Name:        "GoServiceExampleLogging",
		DisplayName: "Go Service Example for Logging",
		Description: "This is an example Go service that outputs log messages.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
