package cmd

import (
	"github.com/spf13/cobra"
	"os"

	"github.com/mercierc/pauli/src"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a pauli project.",
	Long: "Create the .pauli folder. This forlder contains:\n" +
		"  - config.yaml: A file that contains all necessary " +
		"information about the build image.\n" +
		"  - pauli.sh: A shell file with predefined commun functions " +
		"to populate.",
	Run: func(cmd *cobra.Command, arg []string) {
		src.InitiateProject(os.Stdin)
	},
}
