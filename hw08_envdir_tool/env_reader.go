package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envMap := make(map[string]EnvValue, 0)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		err = checkFileName(f.Name())
		if err != nil {
			log.Println("checkFileName ", err)
			return nil, err
		}
		if !f.IsDir() {
			pathToEnvFile := path.Join(dir, f.Name())
			fileInfo, err := f.Info()
			if err != nil {
				log.Println("fileInfo ", err)
				return nil, err
			}
			if fileInfo.Size() == 0 {
				envMap[f.Name()] = EnvValue{"", true}
				continue
			}
			content, err := processFileContent(pathToEnvFile)
			if err != nil {
				log.Println("processFileContent", err)
				return nil, err
			}
			envMap[f.Name()] = EnvValue{content, false}
		}

	}
	return envMap, nil
}

func checkFileName(nameOfFile string) error {
	fmt.Println("checkFileName will implement later")
	return nil
}

func processFileContent(pathToEnvFile string) (content string, err error) {
	fd, err := os.Open(pathToEnvFile)
	if err != nil {
		log.Println("Open ", err)
		return "", err
	}
	bufReader := bufio.NewReader(fd)
	// TODO replace ReadString with ReadSlice
	unprocessedEnvVar, err := bufReader.ReadString(byte('\n'))
	if err != nil && err != io.EOF {
		log.Println("bufReader", err)
		return "", err
	}
	trimmedEnVar := strings.TrimRight(unprocessedEnvVar, " \t\n")

	return trimmedEnVar, nil
}
