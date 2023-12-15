package src

import(
	"context"
	"os"
	
	"gopkg.in/yaml.v2"
	"github.com/docker/docker/client"
	// "github.com/docker/docker/errdefs"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/pkg/stdcopy"

	"github.com/mercierc/pauli/logs"
)

type ContainerManager struct {
	cli *client.Client // docker client
	ctx context.Context
	containerID string
	containerName string
	cmd []string
	exist bool
}

type Opt func(*ContainerManager)


// Intialize a container manager based on passed options.
func NewContainerManager(options ...Opt) *ContainerManager {
	c := &ContainerManager{exist: false}
	for _, opt := range options {
		opt(c)
	}
	return c
}


// Pass the command to execute in the build container.
func WithCmd(cmd []string) Opt {
	return func(c *ContainerManager) {
		c.cmd = cmd
	}
}


// Intanciate the docker client and create the docker container based on the
// config.yaml file.
func WithConfigYaml(configYamlPath string, containerName string) Opt {
	return func(c *ContainerManager) {
		c.ctx = context.Background()
		c.containerName = containerName
		// Extract configuration from cofnfig.yaml.
		content, err := os.ReadFile(configYamlPath)
		
		if err != nil {
			logs.Logger.Error().Err(err).Msg("error")
			panic(err)
		}
		
		// Initialize the docker client.
		c.cli, err = client.NewClientWithOpts()

		// Parse yaml config  file.
		var confYaml Configuration
		err = yaml.Unmarshal(content, &confYaml)

		// Create Mounts
		mounts := make([]mount.Mount, len(confYaml.Builder.Volumes) + 1)

		for i, volume := range confYaml.Builder.Volumes {
			mounts[i] = mount.Mount{
				Source: volume.Source,
				Target: volume.Target,
				ReadOnly: false,
				Type: "bind",
			}
			logs.Logger.Info().Msgf("%s mounted to %s with type %s",
				volume.Source, volume.Target, volume.Type)
		}
		
		cwd, _ := os.Getwd()
				
		mounts[len(mounts)-1] = mount.Mount{
			Source: cwd,
			Target:"/app",
			ReadOnly: true,
			Type: "bind",
		}
		logs.Logger.Info().Msgf("%s mounted to %s with type %s",
			cwd, "/app", "bind")

		logs.Logger.Debug().Msgf("Command: %s", c.cmd)
		
		// Convert the client.Config
		conf := container.Config{
			AttachStdin: true,  // Attach the standard input, makes possible user interaction
			AttachStdout: true,  // Attach the standard output
			AttachStderr: false,  // Attach the standard error
			Env: []string{"ENVVAR=VALEUR"},  // List of environment variable to set in the container
			Cmd: c.cmd,  // Command to run when starting the container
			Image:"alpine_pauli:latest",  // Name of the image 
			WorkingDir: "/app",  // Current directory (PWD) in the command will be launched
			Tty: false,
		}
		confHost := container.HostConfig{Mounts: mounts}


		// Create a new valid container
		resp, err := c.cli.ContainerCreate(c.ctx, &conf, &confHost, nil, nil, c.containerName)
		if err != nil {
			
			logs.Logger.Error().Err(err).Msg("")
		}
		
		c.containerID = resp.ID
		
		logs.Logger.Info().Msgf("Container %s created with ID=%s",
			c.containerName,
			c.containerID[:10],
		)
	}
}


func (c *ContainerManager) Start() {
	logs.Logger.Trace().Msgf("Start container %v", c.containerName)
	
	err := c.cli.ContainerStart(c.ctx, c.containerID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}
	
	c.DockerLogsToHost()
	
	// Remove already existing container.
	err = c.cli.ContainerRemove(
		c.ctx,
		c.containerName,
		types.ContainerRemoveOptions{Force: true})
	logs.Logger.Error().Err(err)
	
}


// Extract ID of the container if exists.
func (c *ContainerManager) GetID() string {

	containers,  _ := c.cli.ContainerList(context.Background(),
		types.ContainerListOptions{
			All: true,
		})
	for _, container := range containers {
		// If same container found, remove it.
		if container.Names[0][1:] == c.containerName {
			return container.ID
		}
	}
	logs.Logger.Debug().Msgf("Container %s not found.", c.containerName)
	return ""
}

// Execute a command on an already existing container.
func (c *ContainerManager) Exec(cmd []string) {
	logs.Logger.Trace().Msgf("Exec container %v", c.containerName)

	execConfig := types.ExecConfig{
		Privileged: false,     // Is the container in privileged mode
		Tty: false,     // Attach standard streams to a tty.
		AttachStdin: false,     // Attach the standard input, makes possible user interaction
		AttachStderr: false,     // Attach the standard error
		AttachStdout: false,     // Attach the standard output
		Detach: false,    // Execute in detach mode
		Env: []string{}, // Environment variables
		WorkingDir: "/app",   // Working directory
		Cmd: cmd, // Execution commands and args
	}

	resp, err := c.cli.ContainerExecCreate(c.ctx, c.containerID, execConfig)
	if err == nil {
		panic(err)
	}
	logs.Logger.Info().Msgf("ExecCreate returns: %+v", resp)
	err = c.cli.ContainerExecStart(c.ctx, resp.ID, types.ExecStartCheck{})
	if err == nil {
		panic(err)
	}

	
	ci, err := c.cli.ContainerExecInspect(c.ctx, resp.ID)
	if err == nil {
		panic(err)
	}

	logs.Logger.Info().Msgf("Exec process returns: %+v", ci)

	c.DockerLogsToHost()
}

// Write docker logs on the host terminal.
func (c *ContainerManager) DockerLogsToHost() {
	statusCh, errCh := c.cli.ContainerWait(c.ctx, c.containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := c.cli.ContainerLogs(
		c.ctx,
		c.containerID,
		types.ContainerLogsOptions{
			ShowStdout: true,
			Follow: true,
		})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}
