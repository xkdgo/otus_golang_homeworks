package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const (
	CmdErrorCode = 1
)

var errParseEnviron = errors.New("couldnt parse environs")

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	var exitErrorTarget *exec.ExitError
	if len(cmd) == 0 {
		return CmdErrorCode
	}
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	currentEnvs := os.Environ()
	envMap, err := collectCommandEnvMap(currentEnvs)
	if err != nil {
		log.Println(err)
		return CmdErrorCode
	}
	for envname, envvalue := range env {
		switch envvalue.NeedRemove {
		case true:
			delete(envMap, envname)
		case false:
			envMap[envname] = envvalue.Value
		}
	}
	command.Env = environ(envMap)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Start(); err != nil {
		log.Println(err)
	}
	if err := command.Wait(); err != nil {
		if errors.As(err, &exitErrorTarget) {
			returnCode = exitErrorTarget.ExitCode()
		} else {
			log.Println(err)
		}
	}
	return returnCode
}

func collectCommandEnvMap(envs []string) (map[string]string, error) {
	envsMap := make(map[string]string, len(envs))
	for _, envItem := range envs {
		keyVal := strings.SplitN(envItem, "=", 2)
		if len(keyVal) != 2 {
			return nil, errParseEnviron
		}
		envsMap[keyVal[0]] = keyVal[1]
	}
	return envsMap, nil
}

func environ(envMap map[string]string) []string {
	environSlice := make([]string, 0, len(envMap))
	b := strings.Builder{}
	for envName, envVal := range envMap {
		b.WriteString(envName)
		b.WriteString("=")
		b.WriteString(envVal)
		environSlice = append(environSlice, b.String())
		b.Reset()
	}
	sort.Strings(environSlice)
	return environSlice
}
