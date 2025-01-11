package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/jamesrr39/go-errorsx"
	"github.com/jamesrr39/yarn-contained/docker"
	_ "github.com/joho/godotenv/autoload"
)

// Env variable names
const (
	EnvVarLoggerLevelVerbose      = "YARN_CONTAINED_VERBOSE_LOGGING"
	EnvVarForceDockerImageRebuild = "YARN_CONTAINED_FORCE_DOCKER_BUILD"
	EnvVarDockerTool              = "YARN_CONTAINED_DOCKERTOOL"
	EnvVarDockerPortForward       = "YARN_CONTAINED_PORT_FORWARD"
	EnvVarExtraContainerArgs      = "YARN_CONTAINED_EXTRA_CONTAINER_ARGS"
	EnvVarForwardedEnvVars        = "YARN_CONTAINED_ENV_VARS" // comma separated, e.g. "NPM_TOKEN,AWS_SECRET_KEY"
)

const (
	DockerImageName = "jamesrr39/yarncontained"
	ProjectURL      = "https://github.com/jamesrr39/yarn-contained"
)

var (
	forceDockerRebuild bool
	portForward        string
)

func getLoggerWriter() io.Writer {
	if envBoolean(EnvVarLoggerLevelVerbose) {
		return os.Stderr
	}

	return io.Discard
}

func trimSpaces(items []string) []string {
	var out []string

	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}

	return out
}

func exitIfErrWithoutStackTrace(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func main() {
	log.SetOutput(getLoggerWriter())

	log.Printf("using yarn-contained: %s\n", ProjectURL)

	var err error

	forceDockerRebuild = envBoolean(EnvVarForceDockerImageRebuild)
	portForward = envString(EnvVarDockerPortForward, "")
	extraContainerArgs := trimSpaces(strings.Split(envString(EnvVarExtraContainerArgs, ""), " "))

	dockerTool, err := getDockerTool()
	errorsx.ExitIfErr(errorsx.Wrap(err))

	log.Printf("using %q as the container tool\n", dockerTool)

	dockerService := docker.NewDockerService(dockerTool)

	yarnArgs := os.Args[1:]

	subCommand := ""
	if len(yarnArgs) > 0 {
		subCommand = yarnArgs[0]
	}

	switch subCommand {
	case "init", "create", "--version":
		// continue without package.json, since these commands do not expect a package.json in the working directory.
	default:
		yarnLockExists, err := checkForPackageJson()
		errorsx.ExitIfErr(err)

		if !yarnLockExists {
			fmt.Fprintf(os.Stderr, "%s does not exist in the current working directory and the command was not 'init'. Exiting.\n", packageJsonFilename)
			os.Exit(1)
		}
	}

	err = dockerService.EnsureDockerImage(DockerImageName, forceDockerRebuild)
	errorsx.ExitIfErr(errorsx.Wrap(err))

	workingDir, err := os.Getwd()
	errorsx.ExitIfErr(errorsx.Wrap(err))

	hostUser, err := user.Current()
	errorsx.ExitIfErr(errorsx.Wrap(err))

	err = dockerService.RunImage(DockerImageName, workingDir, yarnArgs, hostUser, portForward, getEnvVarsToForward(), extraContainerArgs)
	exitIfErrWithoutStackTrace(err)
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
)

// look for podman first - if someone has podman installed they probably want to use that.
// then use docker.
var containerApplications = []string{"podman", "docker"}

func getEnvVarsToForward() []docker.EnvironmentVariable {
	envVarNames := strings.Split(os.Getenv(EnvVarForwardedEnvVars), ",")

	var envVars []docker.EnvironmentVariable

	for _, envVarName := range envVarNames {
		envVarName = strings.TrimSpace(envVarName)
		if envVarName == "" {
			continue
		}

		val := os.Getenv(envVarName)

		envVars = append(envVars, docker.EnvironmentVariable{
			Key:   envVarName,
			Value: val,
		})
	}

	return envVars
}

func getDockerTool() (string, errorsx.Error) {
	chosenDockerTool := envString(EnvVarDockerTool, "")
	if chosenDockerTool != "" {
		return chosenDockerTool, nil
	}

	for _, executable := range containerApplications {
		fullPath, err := exec.LookPath(executable)

		if err != nil {
			_, ok := err.(*exec.Error)
			if ok {
				// some error running the executable, maybe it didn't exist, or not enough permissions to run it. Try the next option.
				continue
			}

			return "", errorsx.Wrap(err)
		}

		return fullPath, nil
	}

	return "", errorsx.Errorf("no suitable docker tool found")
}

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
