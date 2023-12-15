package main


import (

	"github.com/mercierc/pauli/cmd"
	"github.com/mercierc/pauli/logs"
)


func main() {
	// CLI
	
	cmd.Parse()
	
	//src.BuildContainerCommand(".pauli/config.yaml")
	logs.Logger.Info().Msg("End")
}


// A faire
// Firnir contaier_manager.go
// On veut creer un conteneur, lancer une commander dans le conteneur et afficher le résultat
// sur le terminal.

// Finir les LogLevel. root.go et logs.go
