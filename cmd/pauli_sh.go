package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/mercierc/pauli/logs"
	"github.com/mercierc/pauli/src"
)

var (
	envVars    []string
	currentCmd string
)

func commonRun(cmd *cobra.Command, args []string) {
	logs.Logger.Debug().Msgf("pauli command %s", args)
	logs.Logger.Debug().Msgf("--env=%s", envVars)

	// Ensure .pauli/pauli.sh and .pauli/config.yml exist.	
	for _, file := range []string{".pauli", ".pauli/config.yaml", ".pauli/pauli.sh"} {
		_, err := os.Stat(file)
		
		if os.IsNotExist(err) {
			logs.Logger.Error().Err(err).Msgf("To solve the issue: " +
				"add the missing file or directory '%v' manually " +
				"or run pauli init to initiate your project", file)
			os.Exit(1)
		}
	}

	var cm *src.ContainerManager

	// Extract container name from the current folder
	containerName, _ := os.Getwd()
	containerName = filepath.Base(containerName) + "_build"

	cm = src.NewContainerManager(
		src.WithName(containerName),
		src.WithEnv(envVars),
		src.WithConfigYaml(configPath, false),
		src.WithCmd(append([]string{"/bin/sh", ".pauli/pauli.sh", currentCmd}, args...)),
	)
	cm.Start()
	cm.Exec()
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Execute the build from pauli.sh",
	Long: "Launch a build container and execute the build function " +
		"from pauli.sh.",
	PreRun: func(cmd *cobra.Command, args []string) {
		currentCmd = "build"
	},
	Run: commonRun,
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute the run from pauli.sh",
	Long: "Launch a build container and execute the run function " +
		"from pauli.sh.",
	PreRun: func(cmd *cobra.Command, args []string) {
		currentCmd = "run"
	},
	Run: commonRun,
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Execute the clean from pauli.sh",
	Long: "Launch a build container and execute the clean function " +
		"from pauli.sh.",
	PreRun: func(cmd *cobra.Command, args []string) {
		currentCmd = "clean"
	},
	Run: commonRun,
}

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Execute the lint from pauli.sh",
	Long: "Launch a build container and execute the lint function " +
		"from pauli.sh.",
	PreRun: func(cmd *cobra.Command, args []string) {
		currentCmd = "lint"
	},
	Run: commonRun,
}

var unittestsCmd = &cobra.Command{
	Use:   "unittests",
	Short: "Execute the unittests from pauli.sh",
	Long: "Launch a build container and execute the unittests function " +
		"from pauli.sh.",
	PreRun: func(cmd *cobra.Command, args []string) {
		currentCmd = "unittests"
	},
	Run: commonRun,
}

var inttestsCmd = &cobra.Command{
	Use:   "inttests",
	Short: "Execute the inttests from pauli.sh",
	Long: "Launch a build container and execute the inttests function " +
		"from pauli.sh.",
	PreRun: func(cmd *cobra.Command, args []string) {
		currentCmd = "inttests"
	},
	Run: commonRun,
}

var staticanalysisCmd = &cobra.Command{
	Use:   "staticanalysis",
	Short: "Execute the staticanalysis from pauli.sh",
	Long: "Launch a build container and execute the staticanalysis function " +
		"from pauli.sh.",
	PreRun: func(cmd *cobra.Command, args []string) {
		currentCmd = "staticanalysis"
	},
	Run: commonRun,
}

var shellCmd = &cobra.Command{
	Use: "shell",
	Short: "Interactive session in the build image that take your " +
		"config.yml file into account.",
	Long: "Launch an interacive shell in the build container as you would " +
		"do it with the -it option. The first argument allows to " +
		"choose between sh and bash. By default sh\n" +
		"Example: pauli shell [sh|bash]",
	Args:      cobra.MaximumNArgs(1),
	ValidArgs: []string{"sh", "bash"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			args = []string{"sh"}
		}
		containerName, _ := os.Getwd()
		containerName = filepath.Base(containerName) + "_build"

		logs.Logger.Trace().Msgf("shell %s", args[0])
		cm := src.NewContainerManager(
			src.WithName(containerName),
			src.WithEnv(envVars),
			src.WithConfigYaml(configPath, true),
		)
		cm.Shell(args[0])
	},
}
