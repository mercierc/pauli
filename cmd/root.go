package cmd

import(
	"fmt"

	"github.com/spf13/cobra"
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
	Version: "0.0.0",
}


// Execute executes the root command.
func Execute() error {
	rootCmd.AddCommand(initCmd)
	return rootCmd.Execute()
}




