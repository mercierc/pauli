package src

import (

	"fmt"
	"os"
	"bufio"
	"io/ioutil"
	"testing"
	"gopkg.in/yaml.v2"

	"github.com/mercierc/pauli/logs"
)



func TestInitiateProject(t *testing.T) {

	// Read inputs from temporary file
	file, err := ioutil.TempFile("/tmp", "test")
	defer file.Close()
	if err != nil {
		logs.Logger.Error().Err(err).Msg("error")	
	}

	// Write and replace the cursor at the begining of the file
	_, err = file.WriteString("War\nAnd\nPeace")
	file.Seek(0, 0)

	InitiateProject(file)

	_, err = os.Stat(".pauli")
	
	if os.IsNotExist(err) {
	        t.Fatalf(".pauli folder does not exist %v", err)
	}

	// Load template from config.yaml
	content, err := os.ReadFile(".pauli/config.yaml")
	// Parse yaml config  file.
	var confYaml Configuration
	err = yaml.Unmarshal(content, &confYaml)

	if (
		confYaml.Name != "War" ||
		confYaml.Builder.Image != "And" ||
		confYaml.Builder.Tag != "Peace") {
		
		t.Fatal("Input for config.yaml are wrong.")
	}
	fmt.Printf("%+v", confYaml)

	// By defaut
	InitiateProject(file)

	// Load template from config.yaml
	content, err = os.ReadFile(".pauli/config.yaml")
	// Parse yaml config  file.
	err = yaml.Unmarshal(content, &confYaml)
	fmt.Printf("%+v", confYaml)
	if (
		confYaml.Name != "src" ||
		confYaml.Builder.Image != "<image_name>" ||
		confYaml.Builder.Tag != "latest") {
		
		t.Fatalf("Input for config.yaml are wrong.")
	}

	// Ensure paui.sh is downloaded.
	f, err := os.Open(".pauli/pauli.sh")
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	firstLine := scanner.Text()
	if firstLine != "#!/bin/sh" && firstLine != "#!/bin/bash" {
		t.Fatal("pauli.sh is empty")
	}
	
	//Clean
	os.RemoveAll(".pauli")
}
