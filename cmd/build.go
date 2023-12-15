package cmd

import(
	"github.com/spf13/cobra"

	"github.com/mercierc/pauli/src"
	"github.com/mercierc/pauli/logs"
)


var buildCmd = &cobra.Command{
	Use: "build",
	Short: "Execute the build from pauli.sh",
	Long: "Launch a build container and execute the build function " +
		"pauli.sh from it.",
	Run: func(cmd *cobra.Command, args []string){
		logs.Logger.Trace().Msgf("build %s", args)
		cm := src.NewContainerManager(
			src.WithCmd(append([]string{".pauli/pauli.sh", "build"}, args...)),
			src.WithConfigYaml(configPath, "pauli_build"),
		)

		cm.Start()
	},
}
