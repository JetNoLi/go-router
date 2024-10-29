package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const DEBUG = true

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

// The static content has the {module_name} placholder, this function is
// used to replace all ocuurences of the placeholder with the module name
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

func execOrExit(cmdStr string, dir string) string {
	cmds := strings.Split(cmdStr, " ")

	cmd := exec.Command(cmds[0], cmds[1:]...)

	if dir != "" {
		cmd.Dir = dir
	}

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatalf("Error running command %s: %v\nOutput: %s", cmdStr, err, output)
	}

	if DEBUG {
		fmt.Printf("cmd: %s\n", cmdStr)
		fmt.Printf("output: %s\n", output)
	}

	return string(output)
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("format is: grc create module_name")
		return
	}

	moduleName := os.Args[2]
	projectName := strings.Split(moduleName, "/")[len(strings.Split(moduleName, "/"))-1]

	//TODO: ADD In CLI Flow
	goRouterVersion := "4d34c2583bfbade712b4ec5953318cfaf6378b92"

	err := os.Mkdir(projectName, os.ModePerm)

	if err != nil {
		log.Fatalf("Error running command: %v", err)
	}

	// cmd := exec.Command("go", "mod", "init", moduleName)
	// cmd.Dir = projectName

	cmd := fmt.Sprintf("go mod init %s", moduleName)
	execOrExit(cmd, projectName)

	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	log.Fatalf("Error running command: %v\nOutput: %s", err, output)
	// }

	cmd = fmt.Sprintf("cp -r grc/static/ %s", projectName)
	execOrExit(cmd, "")

	// cmd = exec.Command("cp", "-r", "grc/static/", projectName)

	// output, err = cmd.CombinedOutput()
	// if err != nil {
	// 	log.Fatalf("Error running command: %v\nOutput: %s", err, output)
	// }

	execOrExit("go get github.com/a-h/templ", projectName)

	// cmd = exec.Command("go", "get", "github.com/a-h/templ")
	// cmd.Dir = projectName

	// output, err = cmd.CombinedOutput()
	// if err != nil {
	// 	log.Fatalf("Error running command: %v\nOutput: %s", err, output)
	// }

	cmd = "go get github.com/jetnoli/go-router"

	if goRouterVersion != "" {
		cmd += fmt.Sprintf("@%s", goRouterVersion)
	}

	execOrExit(cmd, projectName)
	// cmd = exec.Command("go", "get", "github.com/jetnoli/go-router@4d34c2583bfbade712b4ec5953318cfaf6378b92")
	// cmd.Dir = projectName

	// output, err = cmd.CombinedOutput()
	// if err != nil {
	// 	log.Fatalf("Error running command: %v\nOutput: %s", err, output)
	// }

	err = replaceModuleName(projectName, moduleName, projectName)

	if err != nil {
		log.Fatalf("Error replacing module name %s, %s", projectName, moduleName)
	}

	// cmd = exec.Command("templ", "generate")
	// cmd.Dir = projectName

	execOrExit("templ generate", projectName)

	// output, err = cmd.CombinedOutput()
	// if err != nil {
	// 	log.Fatalf("Error running command: %v\nOutput: %s", err, output)
	// 	os.Exit(1)
	// }

	execOrExit("go mod tidy", projectName)
	// cmd.Dir = projectName

	// output, err = cmd.CombinedOutput()
	// if err != nil {
	// 	log.Fatalf("Error running command: %v\nOutput: %s", err, output)
	// 	os.Exit(1)
	// }

	fmt.Println("Project Created Successfully!")
}
