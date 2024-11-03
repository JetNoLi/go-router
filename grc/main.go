package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const FixCategory = "addembedimport" // recognized by gopls ApplyFix

// Prevent go pls from removing embed script
var _ = embed.FS{}

//go:embed create-repo.sh
var script []byte

type CmdHandler = func()

const DEBUG = true

var CMDS = map[string]CmdHandler{
	"create": createProject,
}

var BASE_ENV_VARS = map[string]string{
	"PORT":              "3000",
	"TEMPL_VERSION":     "latest",
	"GO_ROUTER_VERSION": "latest",
	"APP_NAME":          "myapp",
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

		fmt.Println("replacing ", fileName)

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

func createEnvFile(fileName string, vars map[string]string) error {
	data := ""

	for key, v := range vars {
		if data != "" {
			data += "\n"
		}

		data += fmt.Sprintf("%s=%s", key, v)
	}

	//TODO: Double check which perms to use
	return os.WriteFile(fileName, []byte(data), os.ModePerm)
}

// This function is used to execute commands with os.Exec
// and return the output or call os.Exit(1) on failure
func execOrExit(cmdStr string, dir string) string {
	cmds := strings.Split(strings.TrimSpace(cmdStr), " ")

	fmt.Println(cmds)

	cmd := exec.Command(cmds[0], cmds[1:]...)

	if dir != "" {
		cmd.Dir = dir
	}

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatalf("Error running command %s in %s %v\nOutput: %s", cmdStr, dir, err, output)
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

	cmd := fmt.Sprintf("bash create-repo.sh %s %s", projectName, moduleName)

	execOrExit(cmd, "")

	err := replaceModuleName(projectName, moduleName, projectName)

	if err != nil {
		log.Fatalf("error replacing module name: %s %s, %s", err.Error(), projectName, moduleName)
	}

	cmd = "go get github.com/jetnoli/go-router"

	if *GRCV != "" {
		cmd += fmt.Sprintf("@%s", *GRCV)
	}

	execOrExit(cmd, projectName)

	vars := BASE_ENV_VARS

	vars["APP_NAME"] = projectName

	if *GRCV != "" {
		vars["GO_ROUTER_VERSION"] = *GRCV
	}

	envPath := fmt.Sprintf("%s/.env", projectName)
	err = createEnvFile(envPath, vars)

	if err != nil {
		log.Fatalf("error creating env file %s", err.Error())
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

	scriptPath := "./create-repo.sh"

	// Write the embedded script to a temporary file
	err := os.WriteFile(scriptPath, script, 0755) // Give it executable permissions
	if err != nil {
		log.Fatalf("Failed to write script to file: %v", err)
	}
	defer os.Remove(scriptPath) // Clean up the temporary file after execution

	process()

	fmt.Println("Project Created Successfully!")
}
