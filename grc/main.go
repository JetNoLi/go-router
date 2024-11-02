package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type CmdHandler = func()

const DEBUG = true

var CMDS = map[string]CmdHandler{
	"create": createProject,
}

// Flags

// go router version commit version to be used
var GRCV = flag.String("cv", "", "commit version to use for go router installation")

// The static content has the incorrect import paths, this function is
// used to replace all occurences with the module name
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

		updatedContent := strings.ReplaceAll(string(content), "github.com/jetnoli/go-router/grc/static", moduleName)

		os.WriteFile(fileName, []byte(updatedContent), os.ModePerm)
	}

	return nil
}

// This function is used to execute commands with os.Exec
// and return the output or call os.Exit(1) on failure
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

func createProject() {
	moduleName := flag.Arg(1)

	if moduleName == "" {
		log.Fatal("no module name provided")
	}

	projectName := strings.Split(moduleName, "/")[len(strings.Split(moduleName, "/"))-1]

	fmt.Println(moduleName, projectName, flag.Args())

	err := os.Mkdir(projectName, os.ModePerm)

	if err != nil {
		log.Fatalf("error running command: %v", err)
	}

	cmd := fmt.Sprintf("go mod init %s", moduleName)
	execOrExit(cmd, projectName)

	cmd = fmt.Sprintf("cp -r grc/static/ %s", projectName)
	execOrExit(cmd, "")

	execOrExit("go get github.com/a-h/templ", projectName)

	cmd = "go get github.com/jetnoli/go-router"

	if *GRCV != "" {
		cmd += fmt.Sprintf("@%s", *GRCV)
	}

	execOrExit(cmd, projectName)

	err = replaceModuleName(projectName, moduleName, projectName)

	if err != nil {
		log.Fatalf("error replacing module name %s, %s", projectName, moduleName)
	}

	execOrExit("templ generate", projectName)

	execOrExit("go mod tidy", projectName)
}

func main() {
	flag.Parse()
	cmd := flag.Arg(0)

	if cmd == "" {
		log.Fatal("no cmd specified")
	}

	process, ok := CMDS[cmd]

	if !ok {
		fmt.Printf("invalid arg %s, allowed commands include:\n", cmd)

		for key := range CMDS {
			fmt.Printf("- %s\n", key)
		}

		os.Exit(1)
	}

	process()

	fmt.Println("Project Created Successfully!")
}
