package docker

import (
	"bytes"
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"

	"github.com/jamesrr39/goutil/errorsx"
)

var (
	//go:embed Dockerfile
	dockerFile string
)

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
		// cmd.Stdout = os.Stdout
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

const containerWorkingDir = "/opt/yarn-contained/workspace"

func (ds *DockerService) RunImage(imageName, workingDir string, yarnArgs []string, hostUID string) errorsx.Error {

	dockerArgs := []string{
		"run",
		"--rm",
		"-v", fmt.Sprintf("%s:%s", workingDir, containerWorkingDir),
		"-e", fmt.Sprintf("USERNAME=user%s", hostUID),
		"-e", fmt.Sprintf("HOST_USER_ID=%s", hostUID),
		imageName,
		"yarn",
	}
	for _, yarnArg := range yarnArgs {
		dockerArgs = append(dockerArgs, yarnArg)
	}

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
