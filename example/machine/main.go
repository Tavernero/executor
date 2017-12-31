package main

import (
	"github.com/m49n3t0/executor/machine"
	log "github.com/sirupsen/logrus"
)

// Main function
func main() {

	log.Info("Begin")

    machine.RunDefault("toto_function")

	log.Info("Done")
}
