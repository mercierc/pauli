package cmd

import(
	"github.com/spf13/cobra"

	"github.com/mercierc/pauli/src"
	"github.com/mercierc/pauli/logs"
)


var shellCmd = &cobra.Command{
	Use: "shell",
	Short: "Interactive session in the build image.",
	Long: "Launch an interacive shell the build container as you would " +
		"do it with the -it option.",
	Run: func(cmd *cobra.Command, args []string){
		logs.Logger.Trace().Msgf("shell %s", args)
		cm := src.NewContainerManager(
			src.WithEntryPoint([]string{"sh", "-c", "'while true; do sleep 10; done'"}),
			src.WithConfigYaml(configPath, "pauli_build", true),
		)
		cm.Shell()
	},
}
