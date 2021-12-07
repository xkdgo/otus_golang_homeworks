package main

import (
	"log"
	"os"
)

func main() {
	// Place your code here.
	argsWithoutScriptName := os.Args[1:]
	// fmt.Println(len(argsWithoutScriptName))
	if len(argsWithoutScriptName) < 2 {
		log.Fatalf("use:\ngo-env directory_with_envs cmd args")
	}
	envDir := argsWithoutScriptName[0]
	cmd := argsWithoutScriptName[1:]
	env, err := ReadDir(envDir)
	if err != nil {
		log.Fatalf("fail to read directory \"%v\". Error is: %v", envDir, err)
	}
	os.Exit(RunCmd(cmd, env))
}
