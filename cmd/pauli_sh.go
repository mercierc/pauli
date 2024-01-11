package cmd

import(
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/mercierc/pauli/src"
	"github.com/mercierc/pauli/logs"
)

var currentCmd string

func commonRun(cmd *cobra.Command, args []string) {
	logs.Logger.Debug().Msgf("pauli command %s", args)
	logs.Logger.Debug().Msgf("--env=%s", envVars)

	var cm *src.ContainerManager
	
	// Extract container name from the currebt folder
	containerName, _ := os.Getwd()
	containerName = filepath.Base(containerName) + "_xxXx"
	
	ID := src.GetID(containerName)
	logs.Logger.Trace().Msgf("ID %v", ID)	
	if ID == "" {
		logs.Logger.Info().Msgf("Create container %v from config.yaml", containerName)	
		cm = src.NewContainerManager(
			src.WithName(containerName),
			src.WithEnv(envVars),
			src.WithConfigYaml(configPath, false),
			src.WithCmd(append([]string{"/bin/sh", ".pauli/pauli.sh", currentCmd}, args...)),
		)
	} else {
		logs.Logger.Info().Msgf("Use already existing %s container", containerName)	
		cm = src.NewContainerManager(
			src.WithName(containerName),
			src.WithEnv(envVars),
			src.WithCmd(append([]string{"/bin/sh", ".pauli/pauli.sh", currentCmd}, args...)),
		)
	}
	cm.Start()
	cm.Exec()
}


var buildCmd = &cobra.Command{
	Use: "build",
	Short: "Execute the build from pauli.sh",
	Long: "Launch a build container and execute the build function " +
	      "from pauli.sh.",
	PreRun: func (cmd *cobra.Command, args []string) {
		currentCmd = "build"
	},
	Run: commonRun,
}


var runCmd = &cobra.Command{
	Use: "run",
	Short: "Execute the run from pauli.sh",
	Long: "Launch a build container and execute the run function " +
		"from pauli.sh.",
	PreRun: func (cmd *cobra.Command, args []string) {
		currentCmd = "run"
	},
	Run: commonRun,
}


var cleanCmd = &cobra.Command{
	Use: "clean",
	Short: "Execute the clean from pauli.sh",
	Long: "Launch a build container and execute the clean function " +
		"from pauli.sh.",
	PreRun: func (cmd *cobra.Command, args []string) {
		currentCmd = "clean"
	},
	Run: commonRun,
}


var lintCmd = &cobra.Command{
	Use: "lint",
	Short: "Execute the lint from pauli.sh",
	Long: "Launch a build container and execute the lint function " +
		"from pauli.sh.",
	PreRun: func (cmd *cobra.Command, args []string) {
		currentCmd = "lint"
	},
	Run: commonRun,
}


var unittestsCmd = &cobra.Command{
	Use: "unittests",
	Short: "Execute the unittests from pauli.sh",
	Long: "Launch a build container and execute the unittests function " +
		"from pauli.sh.",
	PreRun: func (cmd *cobra.Command, args []string) {
		currentCmd = "unittests"
	},
	Run: commonRun,
}


var inttestsCmd = &cobra.Command{
	Use: "inttests",
	Short: "Execute the inttests from pauli.sh",
	Long: "Launch a build container and execute the inttests function " +
		"from pauli.sh.",
	PreRun: func (cmd *cobra.Command, args []string) {
		currentCmd = "inttests"
	},
	Run: commonRun,
}


var staticanalysisCmd = &cobra.Command{
	Use: "staticanalysis",
	Short: "Execute the staticanalysis from pauli.sh",
	Long: "Launch a build container and execute the staticanalysis function " +
		"from pauli.sh.",
	PreRun: func (cmd *cobra.Command, args []string) {
		currentCmd = "staticanalysis"
	},
	Run: commonRun,
}


var(
	envVars []string
)


func init() {
	for _, c := range []*cobra.Command{
		buildCmd,
		runCmd,
		cleanCmd,
		lintCmd,
		unittestsCmd,
		inttestsCmd,
		staticanalysisCmd,
	} {
		c.Flags().StringArrayVarP(&envVars, "env",
			"e", []string{}, "--env K11=V1 --env K2=V2")
		
		rootCmd.AddCommand(c)		
	}
}
