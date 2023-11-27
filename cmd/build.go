package cmd

import(
	"fmt"

	"github.com/spf13/cobra"

	//"github.com/mercierc/pauli/src"
)


var buildCmd = &cobra.Command{
	Use: "build",
	Short: "Execute the build from pauli.sh",
	Long: "Launch a build container and execute the build function " +
		"pauli.sh from it.",
	Run: func(cmd *cobra.Command, arg []string){
		fmt.Println("pauli build, not implemented yet.")
	},
}
