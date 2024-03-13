package src

import (
	"bufio"
	"fmt"
	"github.com/mercierc/pauli/logs"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"testing"
)

// This test the most important Opt. It creates a container from the config.yaml
// file and test that the created container is well configured.
func TestNewContainerManager(t *testing.T) {

	logs.Init("trace", true)
	config := Configuration{
		Builder: Builder{
			Image:      `hello-world`,
			Tag:        `latest`,
			Privileged: true,
			Volumes: []Volume{
				Volume{
					Type:   `bind`,
					Source: `/var/run/docker.sock`,
					Target: `/var/run/docker.sock`,
				},
			},
		},
		Name: `Test_pauli_image`,
	}

	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		t.Fatalf("Error at yaml writing: %v", err)
	}

	file, err := ioutil.TempFile("/tmp", "config.yaml")
	writer := bufio.NewWriter(file)

	_, err = writer.WriteString(string(yamlData))
	writer.Flush()

	defer file.Close()
	defer os.Remove(file.Name())

	cm := NewContainerManager(
		WithName("pauli_TI"),
		WithEnv([]string{"VAR1=1", "VAR2=2"}),
		WithConfigYaml(file.Name(), false),
		WithCmd([]string{"/bin/sh", ".pauli/pauli.sh"}),
	)

	// Here we load the information of the created container and ensure
	// that the config is correct
	json, err := cm.cli.ContainerInspect(cm.ctx, cm.containerName)

	// Check mounted volumes
	cwd, err := os.Getwd()

	// Current working dir is mounted by default.
	mountPoints := []map[string]string{
		map[string]string{
			`Source`:      `/var/run/docker.sock`,
			`Destination`: `/var/run/docker.sock`,
		},
		map[string]string{
			`Source`:      cwd,
			`Destination`: `/app`,
		},
	}

	var mountChecked uint8 = 0
	for _, el := range mountPoints {
		for _, el2 := range json.Mounts {
			if el2.Destination == el["Destination"] && el2.Source == el["Source"] {
				mountChecked += 1
			}
		}
	}
	fmt.Println("mountedCheck", mountChecked)
	if mountChecked != uint8(len(mountPoints)) {
		t.Fatalf("All volumes are not mounted %d/%d", mountChecked, len(mountPoints))
	}

	// Check env variables
	if json.Config.Env[0] != "VAR1=1" || json.Config.Env[1] != "VAR2=2" {
		t.Fatalf("Wrong env variables %v, %v. Waited: %v %v", json.Config.Env[0], json.Config.Env[1], "VAR1=1", "VAR2=2")

	}

	// Check the Cmd
	if json.Config.Cmd[0] != "sleep" || json.Config.Cmd[1] != "infinity" {
		t.Fatalf("Wrong command %v. Waited `sleep infinity`", json.Config.Cmd)
	}

	fmt.Println("Cmd ", json.Config.Cmd)

}
