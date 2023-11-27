package src

import(
	"os/exec"
	"os"
	//"errors"
	"fmt"
	
	"gopkg.in/yaml.v2"
	"github.com/rs/zerolog/log"
)


// Start an interactive session in a docker container from the build image.
func BuildContainerShell(configPath string) error {
	// Read and parse the config.yaml file.
	content, err := os.ReadFile(configPath)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	dockerCmdTmpl := "docker run -t {{ .Volumes }}{{ .Image }}"
	var conf Configuration
	err = yaml.Unmarshal(content, &conf)
	fmt.Printf("Configuration %+v\n", conf.Builder)

	var dockerCommand := "docker run"
	var volumes string
	for _, el := range conf.Builder.Volumes {
		dockerCommand += " -v " + el.Source + ":" + el.Target
	}

	docker
	fmt.Println("Volumes: ", volumes)
	
	return nil
}

// Execute bash functions presents in .pauli/pauli.sh
func RunPauliShell(command string, args []string) {
// Check the pauli folder is present with its files.
	path, err := exec.LookPath("./.pauli/pauli.sh")
	if err != nil {
		log.Fatal().Msg("installing pauli.sh is in your future " + path)
	}
	log.Info().Msgf("Command %s exists", path)

	
	// Construction de la commande Docker
	shell := exec.Command("./.pauli/pauli.sh", "build")

	// Configuration de la sortie standard pour afficher les résultats de la commande
	shell.Stdout = os.Stdout
	shell.Stderr = os.Stderr
	shell.Stdin = os.Stdin

	// Exécution de la commande Docker
	shell.Run()
}

// Former la commander docke rrun ligne 23
