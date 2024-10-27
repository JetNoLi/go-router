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

func initializeProject(name string) error {
	cmd := exec.Command("go", "mod", "init", name)

	// cmd.Dir = directory // Set the directory where you want to initialize the module

	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("failed to initialize module: %w. Output: %s", err, string(output))
	}

	return nil
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

	if len(os.Args) < 2 {
		fmt.Println("Please provide arguments")
		return
	}

	// Access specific arguments
	moduleName := os.Args[1]
	projectName := strings.Split(moduleName, "/")[len(strings.Split(moduleName, "/"))-1]
	// fmt.Println("First argument:", firstArg)

	// Go mod Init
	err := os.Mkdir(projectName, os.ModePerm)

	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = projectName // Set the directory where you want to initialize the module

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error running command: %v\nOutput: %s", err, output)
		os.Exit(1)
	}

	cmd = exec.Command("cp", "-r", "cli/static/", projectName)
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

	cmd = exec.Command("go", "get", "github.com/jetnoli/go-router@98df071dba70df24c762606e24cff79bce8bca87")
	cmd.Dir = projectName

	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error running command: %v\nOutput: %s", err, output)
		os.Exit(1)
	}

	err = replaceModuleName(projectName, moduleName, projectName)

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
