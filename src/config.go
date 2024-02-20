package src

import(
	"os"
	"fmt"
	"text/template"
	"bufio"
	"io"
	"path/filepath"
	"net/http"
	"github.com/mercierc/pauli/logs"
)

type Volume struct { 
	Type   string `yaml:"type"`
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

type Builder struct {
	Image   string `yaml:"image"`
	Tag     string `yaml:"tag"`
	Privileged bool   `yaml:"privileged"`
	Volumes []Volume `yaml:"volumes"`
}
	
type Configuration struct {
	Builder Builder `yaml:"builder"`
	Name string `yaml:"name"`
}

var(
	templateContent =`builder:
  image: {{ if .BuildImage }}{{ .BuildImage }}{{ else }}<image_name>{{ end }}
  tag: {{ if .Tag }}{{ .Tag }}{{ else }}latest{{ end }}
  privileged: true
  volumes:
    - type: bind
      source: /var/run/docker.sock
      target: /var/run/docker.sock
name: {{ .ProjectName }}`
)

// Create the .pauli folder with config.yaml and download pauli.sh
// reader: Allow to read from different inputs.

func InitiateProject(reader io.Reader) error {
	templateContent :=`builder:
  image: {{ if .BuildImage }}{{ .BuildImage }}{{ else }}<image_name>{{ end }}
  tag: {{ if .Tag }}{{ .Tag }}{{ else }}latest{{ end }}
  privileged: true
  volumes:
    - type: bind
      source: /var/run/docker.sock
      target: /var/run/docker.sock
name: {{ .ProjectName }}`

	type Initiate struct {
		ProjectName, BuildImage, Tag string
	}

	// Create the .pauli folder.
	if err := os.Mkdir(".pauli", os.ModePerm); err != nil {
		logs.Logger.Error().Err(err).Msg("error")
	}
	// Create the pauli.sh file.
	file, err := os.Create(".pauli/pauli.sh")
	if err != nil {
		logs.Logger.Error().Err(err).Msg("error")
	}
	defer file.Close()

	// Download the raw pauli.sh file.
	downloaded := make(chan bool)
	go func() {
		r, err := http.Get("https://github.com/mercierc/pauli/raw/main/data/pauli.sh")
		if err != nil {
			logs.Logger.Error().Err(err).Msg("error")
		}

		// Write the body in the pauli.sh file	
		if _, err := io.Copy(file, r.Body); err != nil {
			logs.Logger.Error().Err(err).Msg("error")
		}
		// Ensure the http body is close after the end of the function.
		defer r.Body.Close()

		downloaded <- true
		logs.Logger.Info().Msg("pauli.sh downloaded")
	}()

	i := Initiate{}


	scanner := bufio.NewScanner(reader)
	fmt.Printf("Project name (optional, cwd): ")
	scanner.Scan()
	i.ProjectName = scanner.Text()

	fmt.Printf("Name of the build image (optional, <image_name>): ")
	scanner.Scan()
	i.BuildImage = scanner.Text()

	fmt.Printf("tag (optional, latest): ")
	scanner.Scan()
	i.Tag = scanner.Text()

	if i.ProjectName == "" {
		i.ProjectName, _ = os.Getwd()
		i.ProjectName = filepath.Base(i.ProjectName)
	}

	// Fill the config.tmpl
	tmpl, err := template.New("config.tmpl").Parse(templateContent)
	if err != nil {
		panic(err)
	}
	
	// Apply user entries to the template and save.
	outputFile, err := os.Create(".pauli/config.yaml")
	err = tmpl.Execute(outputFile, i)
	if err != nil {
		panic(err)
	}
	logs.Logger.Info().Msgf("pauli.sh downloaded: %t", <-downloaded)
	return nil
}
