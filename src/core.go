package src

import(
	"os/exec"
	"os"
	"errors"
	"strings"
	"fmt"
	
	"gopkg.in/yaml.v2"
	"github.com/mercierc/pauli/logs"
)


// Start an interactive session in a docker container from the build image.
func BuildContainerCommand(configPath string) (string, error) {

	var dockerCommand string
	
	// Read and parse the config.yaml file.
	content, err := os.ReadFile(configPath)
	if err != nil {
		logs.Logger.Error().Err(err).Msg("error")
		return dockerCommand, err
	}

	var conf Configuration
	err = yaml.Unmarshal(content, &conf)

	// Verify mandatory fields.
	if conf.Builder.Image == "" ||
	   conf.Builder.Tag == "" {
		err = errors.New("A field is missing in " + configPath +
			" among image, tag, volumes")
		logs.Logger.Error().Err(err).Msg("error")
		return dockerCommand, err
	}

	dockerCommand = "docker run --workdir /app "

	// Add the cwd and map it to /app as container workdir.
	conf.Builder.Volumes = append(conf.Builder.Volumes, Volume{
		Type: "bind",
		Source: "$PWD",
		Target: "/app",		
	})
	
	for _, el := range conf.Builder.Volumes {
		dockerCommand += "-v " + el.Source + ":" + el.Target + " "
	}
	if !strings.Contains(dockerCommand, "/var/run/docker.sock") {
		err = errors.New("/var/run/docker.sock not mounted")
		logs.Logger.Error().Err(err).Msg("Add /var/run/docker.sock:/var/run/docker.sock to volumes in " + configPath)
	}
	// Add image and tag.
	dockerCommand += fmt.Sprintf("%s:%s ", conf.Builder.Image, conf.Builder.Tag)

	return dockerCommand, nil
}

// Execute bash functions presents in .pauli/pauli.sh
func PauliShell(command string) {
	
	// // Check the pauli folder is present with its files.
	// path, err := exec.LookPath("./.pauli/pauli.sh")
	// if err != nil {
	// 	logs.Logger.Fatal().Msg("installing pauli.sh is in your future " + path)
	// }
	// logs.Logger.Info().Msgf("Command %s exists", path)

	fmt.Println("Command: ", command)
	// Construction de la commande Docker
	shell := exec.Command(command)
	logs.Logger.Info().Msgf("CI")
	// Configuration de la sortie standard pour afficher les résultats de la commande
	shell.Stdout = os.Stdout
	shell.Stderr = os.Stderr
	shell.Stdin = os.Stdin

	// Exécution de la commande Docker
	out, _ := shell.Output()
	fmt.Println(out)
}

// Former la commander docke rrun ligne 23


