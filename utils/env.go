package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

var once sync.Once

func ReadEnv() {
	once.Do(func() {

		path := os.Getenv("ENV_PATH")

		if path == "" {
			path = ".env"
		}

		rawEnvFile, err := os.ReadFile(path)

		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println("Env File Read Successfully")

		envFileLines := strings.Split(string(rawEnvFile), "\n")

		for _, line := range envFileLines {
			// Ignore Comments
			if line[0] == '#' {
				continue
			}
			lineValues := strings.Split(line, "=")
			key, value := lineValues[0], lineValues[1]
			err = os.Setenv(key, value)

			if err != nil {
				log.Fatal(err.Error())
			}
		}
	})

}
