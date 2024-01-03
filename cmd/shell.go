package cmd

import(
	"github.com/spf13/cobra"

	"github.com/mercierc/pauli/src"
	"github.com/mercierc/pauli/logs"
)


var shellCmd = &cobra.Command{
	Use: "shell",
	Short: "Interactive session in the build image that take your " +
	"config.yml file into account.",
	Long: "Launch an interacive shell the build container as you would " +
		"do it with the -it option. The first argument allows to " +
		"choose between sh and bash. By default sh\n" +
		"Example: pauli shell [sh|bash]",
	Args: cobra.MaximumNArgs(1),
	ValidArgs: []string{"sh", "bash"},
	Run: func(cmd *cobra.Command, args []string){
		if len(args)==0 {
			args = []string{"sh"}
		}
		
		logs.Logger.Trace().Msgf("shell %s", args[0])
		cm := src.NewContainerManager(
			src.WithEntryPoint([]string{"sleep", "infinity"}),
			src.WithConfigYaml(configPath, "pauli_build", true),
		)
		cm.Shell(args[0])
	},
}
