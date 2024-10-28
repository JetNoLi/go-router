package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func CreateStructure() {
	// 	handlers
	//  routes
	// 	assets
	//  css
	//		page.css
	// 	scripts
	//		htmx.js
	//  main.go
	// 	views
	//		pages
	//			home
	//		components
	// 			layout

	// Create Go Project with go mod init
	// Copy over file structure

}

func replaceModuleName(projectName string, moduleName string, path string) error {
	dir, err := os.ReadDir(path + "/")

	if err != nil {
		return err
	}

	for _, file := range dir {
		fileName := path + "/" + file.Name()

		if file.IsDir() {
			err = replaceModuleName(projectName, moduleName, fileName)

			if err != nil {
				return err
			}
			continue
		}

		content, err := os.ReadFile(fileName)

		if err != nil {
			log.Fatalf("error reading file: %s", err)
		}

		updatedContent := strings.ReplaceAll(string(content), "{module_name}", moduleName)

		os.WriteFile(fileName, []byte(updatedContent), os.ModePerm)

	}

	return nil
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("format is: grc create module_name")
		return
	}

	// Access specific arguments
	moduleName := os.Args[2]
	projectName := strings.Split(moduleName, "/")[len(strings.Split(moduleName, "/"))-1]
	// fmt.Println("First argument:", firstArg)

	// Go mod Init
	err := os.Mkdir(projectName, os.ModePerm)

	if err != nil {
		log.Fatalf("Error running command: %v", err)
		os.Exit(1)
	}

	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = projectName // Set the directory where you want to initialize the module

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error running command: %v\nOutput: %s", err, output)
		os.Exit(1)
	}

	cmd = exec.Command("cp", "-r", "grc/static/", projectName)
	// cmd.Dir = projectName

	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error running command: %v\nOutput: %s", err, output)
		os.Exit(1)
	}

	// cmd = exec.Command("go", "get", "github.com/a-h/templ")
	// cmd.Dir = projectName

	// output, err = cmd.CombinedOutput()
	// if err != nil {
	// 	log.Fatalf("Error running command: %v\nOutput: %s", err, output)
	// 	os.Exit(1)
	// }

	cmd = exec.Command("go", "get", "github.com/jetnoli/go-router@4d34c2583bfbade712b4ec5953318cfaf6378b92")
	cmd.Dir = projectName

	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error running command: %v\nOutput: %s", err, output)
		os.Exit(1)
	}

	err = replaceModuleName(projectName, moduleName, projectName)

	if err != nil {
		log.Fatalf("Error replacing module name %s, %s", projectName, moduleName)
		os.Exit(1)
	}

	fmt.Println("replaced moduels")

	cmd = exec.Command("templ", "generate")
	cmd.Dir = projectName

	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error running command: %v\nOutput: %s", err, output)
		os.Exit(1)
	}

	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = projectName

	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error running command: %v\nOutput: %s", err, output)
		os.Exit(1)
	}

}
