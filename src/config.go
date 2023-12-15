package src

import(
	"os"
	"fmt"
	"text/template"
	"io"
	"path/filepath"
	"net/http"

	//"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Volume struct { 
	Type   string `yaml:"type"`
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

type Configuration struct {
	Builder struct {
		Image   string `yaml:"image"`
		Tag     string `yaml:"tag"`
		UseSudo bool   `yaml:"use_sudo"`
		Volumes []Volume `yaml:"volumes"`
	} `yaml:"builder"`
	Name string `yaml:"name"`
}


func InitiateProject() error {

	templateContent :=`builder:
  image: {{ if .BuildImage }}{{ .BuildImage }}{{ else }}<image_name>{{ end }}
  use_sudo: true
  tag: {{ if .Tag }}{{ .Tag }}{{ else }}latest{{ end }}
  volumes:
    - type: bind
      source: /var/run/docker.sock
      target: /var/run/docker.sock
    - type: bind
      source: /root/.docker/config.json
      target: $HOME/.docker/config.json
name: {{ .ProjectName }}`

	type Initiate struct {
		ProjectName, BuildImage, Tag string
	}

	// Create the .pauli folder.
	if err := os.Mkdir(".pauli", os.ModePerm); err != nil {
		log.Error().Err(err)
	}
	// Create the pauli.sh file.
	file, err := os.Create(".pauli/pauli.sh")
	if err != nil {
		log.Error().Err(err)
	}
	defer file.Close()
	
	// Download the raw pauli.sh file.
	var resp *http.Response
	go func(r *http.Response) {

		r, err := http.Get("https://github.com/mercierc/pauli/raw/main/data/pauli.sh")
		if err != nil {
			log.Error().Err(err)
		}

		// Write the body in the pauli.sh file	
		if _, err := io.Copy(file, r.Body); err != nil {
			log.Error().Err(err)
		}
		// Ensure the http body is close after the end of the function.
		defer r.Body.Close()	

	}(resp)

	i := Initiate{}
	fmt.Printf("Project name (optional, cwd): ")
	fmt.Scanln(&i.ProjectName)
	fmt.Printf("Name of the build image (optional, <image_name>): ")
	fmt.Scanln(&i.BuildImage)
	fmt.Printf("tag (optional, latest): ")
	fmt.Scanln(&i.Tag)
	if i.ProjectName == "" {
		i.ProjectName, _ = os.Getwd()
		i.ProjectName = filepath.Base(i.ProjectName)
	}
	fmt.Printf("%+v", i)
	
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

	return nil
}
