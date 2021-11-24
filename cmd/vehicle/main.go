package main

import "github.com/VJftw/vehicle/internal/logging"

var log = logging.Logger

func main() {
	log.Info().Msg("Hello World")
}
