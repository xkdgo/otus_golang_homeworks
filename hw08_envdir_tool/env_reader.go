package main

import (
	"bufio"
	"bytes"
	"errors"
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
	envMap := make(map[string]EnvValue)
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Println("ReadDir ", err)
		return nil, err
	}
	for _, f := range files {
		if !isCorrectFileName(f.Name()) {
			// ignore bad filenames '=' in filename
			continue
		}
		// ignore Dir in Dir
		if !f.IsDir() {
			pathToEnvFile := path.Join(dir, f.Name())
			fileInfo, err := f.Info()
			if err != nil {
				log.Println("fileInfo ", err)
				return nil, err
			}
			if fileInfo.Size() == 0 {
				envMap[f.Name()] = EnvValue{Value: "", NeedRemove: true}
				continue
			}
			content, err := processFileContent(pathToEnvFile)
			if err != nil {
				log.Println("processFileContent", err)
				return nil, err
			}
			envMap[f.Name()] = EnvValue{Value: content, NeedRemove: false}
		}
	}
	return envMap, nil
}

func isCorrectFileName(nameOfFile string) bool {
	for _, sym := range nameOfFile {
		if sym == '=' {
			return false
		}
	}
	return true
}

func processFileContent(pathToEnvFile string) (string, error) {
	fd, err := os.Open(pathToEnvFile)
	if err != nil {
		log.Println("Open ", err)
		return "", err
	}
	bufReader := bufio.NewReader(fd)
	bytesEnvVar, err := bufReader.ReadBytes(byte('\n'))
	if err != nil && !errors.Is(err, io.EOF) {
		log.Println("bufReader", err)
		return "", err
	}
	bytesEnvVar = bytes.ReplaceAll(bytesEnvVar, []byte{0x00}, []byte{'\n'})
	untrimmedEnvVar := string(bytesEnvVar)
	trimmedEnVar := strings.TrimRight(untrimmedEnvVar, " \t\n")
	return trimmedEnVar, nil
}
