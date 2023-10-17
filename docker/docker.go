package docker

import (
	"bytes"
	_ "embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/jamesrr39/go-errorsx"
)

var (
	//go:embed Dockerfile
	dockerFile string
)

const containerWorkingDir = "/opt/yarn-contained/workspace"

type DockerService struct {
	DockerTool string
}

func NewDockerService(dockerTool string) *DockerService {
	return &DockerService{DockerTool: dockerTool}
}

func (ds *DockerService) EnsureDockerImage(imageName string, forceRebuild bool) errorsx.Error {
	buildDockerImage := forceRebuild

	if !forceRebuild {
		cmd := exec.Command(ds.DockerTool, "image", "ls", imageName, "--format={{.ID}}")
		cmd.Stderr = os.Stderr
		output, err := cmd.Output()
		if err != nil {
			return errorsx.Wrap(err)
		}

		trimmed := strings.TrimSpace(string(output))

		if trimmed == "" {
			// doesn't exist yet
			buildDockerImage = true
		}
	}

	if !buildDockerImage {
		// already created
		return nil
	}

	return ds.CreateDockerImage(imageName)
}

type fileType struct {
	Name, Contents string
	Perm           fs.FileMode
}

func (ds *DockerService) CreateDockerImage(imageName string) errorsx.Error {
	cmd := exec.Command(ds.DockerTool, "build", "-t", imageName, "-")

	cmd.Stdin = bytes.NewBufferString(dockerFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return errorsx.Wrap(err)
	}

	return nil
}

type EnvironmentVariable struct {
	Key, Value string
}

func (ds *DockerService) RunImage(imageName, workingDir string, yarnArgs []string, hostUser *user.User, portForward string, envVars []EnvironmentVariable) errorsx.Error {

	dockerArgs := []string{
		"run",
		"--rm",
		"--interactive",
		"--tty",
		"--userns", "keep-id",
		"-v", fmt.Sprintf("%s:%s", workingDir, containerWorkingDir),
	}
	for _, envVar := range envVars {
		dockerArgs = append(dockerArgs, "-e", fmt.Sprintf("%s=%s", envVar.Key, envVar.Value))
	}

	if portForward != "" {
		dockerArgs = append(dockerArgs, "-p", portForward)
	}
	dockerArgs = append(dockerArgs, imageName, "yarn")
	dockerArgs = append(dockerArgs, yarnArgs...)

	log.Printf("running yarn with: %q\n", strings.Join(append([]string{ds.DockerTool}, dockerArgs...), " "))

	cmd := exec.Command(ds.DockerTool, dockerArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return errorsx.Wrap(err)
	}

	return nil
}
