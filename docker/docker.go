package docker

import (
	"bytes"
	_ "embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/jamesrr39/goutil/errorsx"
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

func (ds *DockerService) RunImage(imageName, workingDir string, yarnArgs []string, hostUID, portForward string) errorsx.Error {

	dockerArgs := []string{
		"run",
		"--rm",
		"--interactive",
		"--tty",
		"-v", fmt.Sprintf("%s:%s", workingDir, containerWorkingDir),
		"-e", fmt.Sprintf("USERNAME=user%s", hostUID),
		"-e", fmt.Sprintf("HOST_USER_ID=%s", hostUID),
	}
	if portForward != "" {
		dockerArgs = append(dockerArgs, "-p", portForward)
	}
	dockerArgs = append(dockerArgs, imageName, "yarn")
	dockerArgs = append(dockerArgs, yarnArgs...)

	log.Printf("dockerArgs: %s\n", strings.Join(dockerArgs, " "))

	cmd := exec.Command(ds.DockerTool, dockerArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	// stdinPipe, err := cmd.StdinPipe()
	// if err != nil {
	// 	return errorsx.Wrap(err)
	// }
	// defer stdinPipe.Close()

	// go func() {
	// 	stdinScanner := bufio.NewScanner(os.Stdin)
	// 	for stdinScanner.Scan() {
	// 		b := stdinScanner.Bytes()

	// 		log.Printf("READ: %s\n", b)
	// 		_, err = stdinPipe.Write(append(b, []byte("\n")...))
	// 		if err != nil {
	// 			panic(errorsx.Wrap(err))
	// 		}
	// 	}

	// 	err = stdinScanner.Err()
	// 	if err != nil {
	// 		panic(errorsx.Wrap(err))
	// 	}
	// }()

	// go func() {
	// 	for {
	// 		time.Sleep(time.Second * 2)
	// 		_, err := stdinPipe.Write([]byte("y\n"))
	// 		if err != nil {
	// 			panic(errorsx.Wrap(err))
	// 		}

	// 		// _, err = cmd.Stdout.Write([]byte("y\n"))
	// 		// if err != nil {
	// 		// 	panic(errorsx.Wrap(err))
	// 		// }
	// 	}
	// }()

	err := cmd.Run()
	if err != nil {
		return errorsx.Wrap(err)
	}

	return nil
}
