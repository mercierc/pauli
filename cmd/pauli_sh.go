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
		"from pauli.sh.",
	Run: func(cmd *cobra.Command, args []string){
		logs.Logger.Debug().Msgf("build %s", args)
		logs.Logger.Debug().Msgf("--env=%s", envVars)
		cm := src.NewContainerManager(
			src.WithEnv(envVars),
			src.WithCmd(append([]string{".pauli/pauli.sh", "build"}, args...)),
			src.WithConfigYaml(configPath, "pauli_build", false),
		)
		cm.Start()
	},
}


var runCmd = &cobra.Command{
	Use: "run",
	Short: "Execute the run from pauli.sh",
	Long: "Launch a build container and execute the run function " +
		"from pauli.sh.",
	Run: func(cmd *cobra.Command, args []string){
		logs.Logger.Debug().Msgf("build %s", args)
		logs.Logger.Debug().Msgf("--env=%s", envVars)
		cm := src.NewContainerManager(
			src.WithEnv(envVars),
			src.WithCmd(append([]string{".pauli/pauli.sh", "run"}, args...)),
			src.WithConfigYaml(configPath, "pauli_build", false),
		)
		cm.Start()
	},
}


var cleanCmd = &cobra.Command{
	Use: "clean",
	Short: "Execute the clean from pauli.sh",
	Long: "Launch a build container and execute the clean function " +
		"from pauli.sh.",
	Run: func(cmd *cobra.Command, args []string){
		logs.Logger.Debug().Msgf("build %s", args)
		logs.Logger.Debug().Msgf("--env=%s", envVars)
		cm := src.NewContainerManager(
			src.WithEnv(envVars),
			src.WithCmd(append([]string{".pauli/pauli.sh", "clean"}, args...)),
			src.WithConfigYaml(configPath, "pauli_build", false),
		)
		cm.Start()
	},
}


var lintCmd = &cobra.Command{
	Use: "lint",
	Short: "Execute the lint from pauli.sh",
	Long: "Launch a build container and execute the lint function " +
		"from pauli.sh.",
	Run: func(cmd *cobra.Command, args []string){
		logs.Logger.Debug().Msgf("build %s", args)
		logs.Logger.Debug().Msgf("--env=%s", envVars)
		cm := src.NewContainerManager(
			src.WithEnv(envVars),
			src.WithCmd(append([]string{".pauli/pauli.sh", "lint"}, args...)),
			src.WithConfigYaml(configPath, "pauli_build", false),
		)
		cm.Start()
	},
}


var unittestsCmd = &cobra.Command{
	Use: "unittests",
	Short: "Execute the unittests from pauli.sh",
	Long: "Launch a build container and execute the unittests function " +
		"from pauli.sh.",
	Run: func(cmd *cobra.Command, args []string){
		logs.Logger.Debug().Msgf("build %s", args)
		logs.Logger.Debug().Msgf("--env=%s", envVars)
		cm := src.NewContainerManager(
			src.WithEnv(envVars),
			src.WithCmd(append([]string{".pauli/pauli.sh", "unittests"}, args...)),
			src.WithConfigYaml(configPath, "pauli_build", false),
		)
		cm.Start()
	},
}


var inttestsCmd = &cobra.Command{
	Use: "inttests",
	Short: "Execute the inttests from pauli.sh",
	Long: "Launch a build container and execute the inttests function " +
		"from pauli.sh.",
	Run: func(cmd *cobra.Command, args []string){
		logs.Logger.Debug().Msgf("build %s", args)
		logs.Logger.Debug().Msgf("--env=%s", envVars)
		cm := src.NewContainerManager(
			src.WithEnv(envVars),
			src.WithCmd(append([]string{".pauli/pauli.sh", "inttests"}, args...)),
			src.WithConfigYaml(configPath, "pauli_build", false),
		)
		cm.Start()
	},
}


var staticanalysisCmd = &cobra.Command{
	Use: "staticanalysis",
	Short: "Execute the staticanalysis from pauli.sh",
	Long: "Launch a build container and execute the staticanalysis function " +
		"from pauli.sh.",
	Run: func(cmd *cobra.Command, args []string){
		logs.Logger.Debug().Msgf("build %s", args)
		logs.Logger.Debug().Msgf("--env=%s", envVars)
		cm := src.NewContainerManager(
			src.WithEnv(envVars),
			src.WithCmd(append([]string{".pauli/pauli.sh", "staticanalysis"}, args...)),
			src.WithConfigYaml(configPath, "pauli_build", false),
		)
		cm.Start()
	},
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
			"e", []string{}, "--env VAR1=VAR1 --env VAR2=VAR2")
		rootCmd.AddCommand(c)
		
	}
}
