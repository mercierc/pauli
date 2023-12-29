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
// Continuer la commande shell :
// go run . build --log error
// go run . shell --log error
