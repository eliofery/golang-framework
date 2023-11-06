package main

import "github.com/eliofery/golang-image/pkg/logging"

func main() {
	log := logging.New("prod")
	log.Debug("Hello, World!")
	log.Info("Hello, World!")
}
