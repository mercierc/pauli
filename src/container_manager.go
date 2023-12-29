package src

import(
	"context"
	"os"
	"os/exec"
	//"fmt"
	
	"gopkg.in/yaml.v2"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
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
	entryPoint []string
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

func WithEntryPoint(entryPoint []string) Opt {
	return func(c *ContainerManager) {
		c.entryPoint = entryPoint
	}
}


// Intanciate the docker client and create the docker container based on the
// config.yaml file.
func WithConfigYaml(configYamlPath string, containerName string, shell bool) Opt {
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
			AttachStdin: false,  // makes possible user interaction
			AttachStdout: false,  // Attach the standard output
			AttachStderr: false,  // Attach the standard error
	                Tty: false,
			Env: []string{"ENVVAR=VALEUR"}, 
			Cmd: c.cmd,  // Command to run when starting the container
			Entrypoint: c.entryPoint,
			Image: confYaml.Builder.Image + ":" + confYaml.Builder.Tag,
			WorkingDir: "/app",  
		}
		confHost := container.HostConfig{Mounts: mounts}


		// Create a new valid container
		resp, err := c.cli.ContainerCreate(c.ctx, &conf, &confHost, nil, nil, c.containerName)

		switch errorType := err.(type) {
		case nil:
			logs.Logger.Info().Msgf("Container %s created with ID=%s",
				c.containerName,
				resp.ID[:10],
			)
		case errdefs.ErrInvalidParameter:
			logs.Logger.Error().Err(err).Msgf("Error type is %T", errorType)
			logs.Logger.Debug().Msg("Adapt your volume mapping configuration to solve this issue")
		default:
			logs.Logger.Error().Err(err).Msgf("Error type is %T", errorType)
			panic(err)
		}
		
		c.containerID = resp.ID
		
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


// Write docker logs on the host terminal.
func (c *ContainerManager) DockerLogsToHost() {
	statusCh, errCh := c.cli.ContainerWait(
		c.ctx,
		c.containerID,
		container.WaitConditionNotRunning,
	)

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
func (c *ContainerManager) Shell() {
	err := c.cli.ContainerStart(c.ctx, c.containerID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("docker", "exec", "-ti", c.containerName, "sh")

	logs.Logger.Info().Msgf("Interactive session >>>")

	// Pipe the standard input/output to the application standar input/output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Exécution de la commande Docker
	err = cmd.Run()

	if err != nil {
		logs.Logger.Error().Err(err)
	}

	// Remove already existing container.
	err = c.cli.ContainerRemove(
		c.ctx,
		c.containerName,
		types.ContainerRemoveOptions{Force: true})
	logs.Logger.Error().Err(err)
}
