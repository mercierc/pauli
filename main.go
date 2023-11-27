package main


import (
	"os"
	
	"github.com/rs/zerolog"

	"github.com/mercierc/pauli/cmd"
	"github.com/mercierc/pauli/src"
)


func main() {

	// CLI
	cmd.Execute()

	environment := "dev"
	var logger zerolog.Logger
	
	if environment == "dev" {
		// Initiate the logger.
		cw := zerolog.ConsoleWriter{
			Out: os.Stderr,
			NoColor: false,
		}
		logger = zerolog.New(cw).Level(zerolog.InfoLevel)

	} else {
		logger = zerolog.New(os.Stderr).Level(zerolog.InfoLevel)
	}

	src.BuildContainerShell(".pauli/config.yaml")
	logger.Info().Msg("End")
}


// A faire
// Aller dans cmd/build.go et src/core pour réléfichir à la fonction d'appel des fonctions bash dans pauli.sh.
// Creer le mécanisme de rentré dans une image en mode intéractif avant d'excuter les pauli.sh
