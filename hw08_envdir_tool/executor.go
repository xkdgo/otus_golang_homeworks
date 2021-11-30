package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

const (
	CmdErrorCode = 1
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	var exitErrorTarget *exec.ExitError
	if len(cmd) == 0 {
		return CmdErrorCode
	}
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	for envname, envvalue := range env {
		switch envvalue.NeedRemove {
		case true:
			os.Unsetenv(envname)
		case false:
			os.Setenv(envname, envvalue.Value)
		}
	}
	command.Env = os.Environ()
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
