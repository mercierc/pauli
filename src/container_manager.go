package src

import(
	"context"
	"os"
	"os/exec"
	"io"
	"time"
	
	"gopkg.in/yaml.v2"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"

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
	env []string
}

type Opt func(*ContainerManager)


// Intialize a container manager based on passed options.
func NewContainerManager(options ...Opt) *ContainerManager {
	c := &ContainerManager{exist: false}

	c.ctx = context.Background()

	// Initialize the docker client.
	c.cli, _ = client.NewClientWithOpts()

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


func WithEnv(env []string) Opt {
	return func(c *ContainerManager) {
		c.env = env
	}
}


func WithName(containerName string) Opt {
	return func(c *ContainerManager) {
		c.containerName = containerName
	}
}


func WithEntryPoint(entryPoint []string) Opt {
	return func(c *ContainerManager) {
		c.entryPoint = entryPoint
	}
}


// Intanciate the docker client and create the docker container based on the
// config.yaml file.
func WithConfigYaml(configYamlPath string, shell bool) Opt {
	return func(c *ContainerManager) {
		// Extract configuration from cofnfig.yaml.
		content, err := os.ReadFile(configYamlPath)
		
		if err != nil {
			logs.Logger.Error().Err(err).Msg("error")
			panic(err)
		}
		
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
			ReadOnly: false,
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
	                Tty: true,
			Env: c.env, 
			Cmd: []string{"sleep", "infinity"},  //  Command to run when starting the container
			Entrypoint: c.entryPoint,
			Image: confYaml.Builder.Image + ":" + confYaml.Builder.Tag,
			WorkingDir: "/app",
		}
		privileged := false
		privileged = privileged || confYaml.Builder.Privileged

		confHost := container.HostConfig{Mounts: mounts, Privileged: privileged}

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
			logs.Logger.Debug().Msg("Adapt your volume mapping configuration to solve this issue.")
		case errdefs.ErrNotFound:
			logs.Logger.Info().Msgf("%s image not present. Pull and retry ", conf.Image)
			logs.Logger.Info().Msgf("Pull %s", conf.Image)
			reader, _ := c.cli.ImagePull(
				c.ctx,
				conf.Image,
				types.ImagePullOptions{},
			)
			io.Copy(os.Stdout, reader)
			resp, err = c.cli.ContainerCreate(c.ctx, &conf, &confHost, nil, nil, c.containerName)
		default:
			logs.Logger.Error().Err(err).Msgf("Error type is %T", errorType)
			panic(err)
		}
		
		c.containerID = resp.ID	
	}
}


func (c *ContainerManager) Start() {
	logs.Logger.Trace().Msgf("Start container %v", c.containerName)	
	
	err := c.cli.ContainerStart(c.ctx, c.containerName, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}
	//c.DockerLogsToHost()

}

func (c *ContainerManager) Exec() {
	logs.Logger.Trace().Msgf("Exec command %v", c.cmd)
	logs.Logger.Trace().Msgf("c.containerID %v", c.containerID)	
	exec, err := c.cli.ContainerExecCreate(
		c.ctx,
		c.containerName,
		types.ExecConfig{
			AttachStdin: false,  // makes possible user interaction
			AttachStdout: true,  // Attach the standard output
			AttachStderr: true,  // Attach the standard error
	                Tty: true,
			Env: c.env, 
			Cmd: c.cmd, //   Command to run when starting the container
			WorkingDir: "/app",
		},
	)

	if err != nil {
		panic(err)
	}	

	hijack, err := c.cli.ContainerExecAttach(c.ctx, exec.ID, types.ExecStartCheck{})	

	err = c.cli.ContainerExecStart(c.ctx, exec.ID, types.ExecStartCheck{Detach: false, Tty: false})

	if err != nil {
		panic(err)
	}

	go func() {
		defer hijack.Conn.Close()
		io.Copy(os.Stdout, hijack.Reader)
	}()
	
	go func() {
		for {
			execInspect, err := c.cli.ContainerExecInspect(c.ctx, exec.ID)
			logs.Logger.Trace().Msgf("Exec 'pauli %v' is running=%v", c.cmd[len(c.cmd)-1], execInspect.Running)	
			if err != nil {
				panic(err)
			}

			if !execInspect.Running {
				break
			}

			time.Sleep(1 * time.Second)
		}
		timeout := 1
		c.cli.ContainerStop(c.ctx, c.containerName, container.StopOptions{ Timeout: &timeout})
		logs.Logger.Info().Msgf("Container %v is stopping", c.containerName)
	}()
	
	c.DockerLogsToHost()
}


// Write docker logs on the host terminal.
func (c *ContainerManager) DockerLogsToHost() {
	out, err := c.cli.ContainerLogs(
		c.ctx,
		c.containerName,
		types.ContainerLogsOptions{
			ShowStdout: true,
			Follow: true,
			Timestamps: true,
			Since: "0s",
		})
	if err != nil {
		panic(err)
	}

	go func() {
		defer out.Close()
		io.Copy(os.Stdout, out)
	}()

	statusCh, errCh := c.cli.ContainerWait(
		c.ctx,
		c.containerName,
		container.WaitConditionNotRunning,
	)
	
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}
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
func (c *ContainerManager) Shell(shell string) {
	err := c.cli.ContainerStart(c.ctx, c.containerID, types.ContainerStartOptions{})
	if err != nil {
		logs.Logger.Error().Err(err).Msgf("Error type is %T", err)
		panic(err)
	}

	cmd := exec.Command("docker", "exec", "--privileged", "-ti", c.containerName, shell)

	logs.Logger.Info().Msgf("Interactive session >>>")

	// Pipe the standard input/output to the application standar input/output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()

	if err != nil {
		logs.Logger.Error().Err(err)
	}

	// Remove container once the interactive session is finished.
	err = c.cli.ContainerRemove(
		c.ctx,
		c.containerName,
		types.ContainerRemoveOptions{Force: true})
	logs.Logger.Error().Err(err)
}


func GetID(containerName string) string {

	cli, _ := client.NewClientWithOpts()
	containers,  _ := cli.ContainerList(context.Background(),
		types.ContainerListOptions{
			All: true,
		})
	for _, container := range containers {
		// If same container found, remove it.
		if container.Names[0][1:] == containerName {
			return container.ID
		}
	}
	logs.Logger.Debug().Msgf("Container %s not found.", containerName)
	return ""
}
