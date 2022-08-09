package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/jamesrr39/goutil/errorsx"
	"github.com/jamesrr39/yarn-contained/docker"
)

const (
	DockerImageName                   = "jamesrr39/yarncontained"
	ForceDockerImageRebuildEnvVarName = "YARN_CONTAINED_FORCE_DOCKER_BUILD"
	DockerToolEnvVarName              = "YARN_CONTAINED_DOCKERTOOL"
	ProjectURL                        = "https://github.com/jamesrr39/yarn-contained"
)

var (
	forceDockerRebuild bool
)

func main() {
	log.Printf("using yarn-contained: %s\n", ProjectURL)

	var err error
	flag.Parse()

	forceDockerRebuild = envBoolean(ForceDockerImageRebuildEnvVarName)

	subCommand := flag.Arg(0)

	dockerService := docker.NewDockerService(envString(DockerToolEnvVarName, dockerTool))

	yarnArgs := os.Args[1:]

	switch subCommand {
	case "init":
		// if len(yarnArgs) != 1 {
		// 	log.Fatalln("sorry, yarn init with extra arguments is currently not supported. Please use just 'yarn init' and then edit package.json by hand")
		// }
		// yarnArgs = append(yarnArgs, "--yes")
	default:
		yarnLockExists, err := checkForPackageJson()
		errorsx.ExitIfErr(err)

		if !yarnLockExists {
			log.Fatalf("%s does not exist in the current working directory and the command was not 'init'. Exiting.\n", packageJsonFilename)
		}
	}

	err = dockerService.EnsureDockerImage(DockerImageName, forceDockerRebuild)
	errorsx.ExitIfErr(errorsx.Wrap(err))

	workingDir, err := os.Getwd()
	errorsx.ExitIfErr(errorsx.Wrap(err))

	hostUser, err := user.Current()
	errorsx.ExitIfErr(errorsx.Wrap(err))

	err = dockerService.RunImage(DockerImageName, workingDir, yarnArgs, hostUser.Uid)
	errorsx.ExitIfErr(errorsx.Wrap(err))
}

func checkForPackageJson() (bool, errorsx.Error) {
	_, err := os.Stat(packageJsonFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, errorsx.Wrap(err)
	}

	return true, nil
}

const (
	packageJsonFilename = "package.json"
	dockerTool          = "docker"
)

func envBoolean(key string) bool {
	value := os.Getenv(key)
	switch strings.ToLower(value) {
	case "1", "true", "on", "yes":
		return true
	default:
		return false
	}
}

func envString(key, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	return value
}
