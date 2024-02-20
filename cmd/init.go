package cmd

import(
	"os"
	"github.com/spf13/cobra"

	"github.com/mercierc/pauli/src"
)


var initCmd = &cobra.Command{
	Use: "init",
	Short: "Initialize a pauli project.",
	Long: "Create the .pauli folder. This forlder contains:\n" +
		"  - config.yaml: A file that contains all necessary " +
		"information about the build image.\n" + 
                "  - pauli.sh: A shell file with predefined commun functions " +
		"to populate.",
	Run: func(cmd *cobra.Command, arg []string){
		src.InitiateProject(os.Stdin)
	},
}
