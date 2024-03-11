package main

import (
	"github.com/mercierc/pauli/cmd"
	"github.com/mercierc/pauli/logs"
)

func main() {
	// CLI
	cmd.Parse()
	logs.Logger.Info().Msg("End")
}
