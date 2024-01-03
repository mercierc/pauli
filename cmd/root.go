package cmd

import(
	"fmt"
	
	"github.com/spf13/cobra"
	"github.com/mercierc/pauli/logs"
)

var rootCmd = &cobra.Command{
	Use: "pauli",
	Short: "Pauli allows you to develop your project in a reproductible" +
		" environment",
	Long: "Pauli allow to transparently develop and run project in a " +
		"docker container containing all the dependencies of your application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Rootcmd")
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize the logger.
		logs.Init(logLevel, dev)
	},
	Version: "0.0.0",
}

var pauliShPath = ".pauli/pauli.sh"
var configPath = ".pauli/config.yaml"

var(
	logLevel string
	dev bool
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&dev,
		"dev", false,
		"Developper friendly log format.")
	
	rootCmd.PersistentFlags().StringVar(&logLevel,
		"log", "info",
		"Niveau de log (trace, debug, info, warn, error, panic)")
}

// Parse the command line.
func Parse() error {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(shellCmd)
	return rootCmd.Execute()
}




