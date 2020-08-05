// +build mage

package main

import (
	. "fmt"
	"os"
	"os/exec"
	"path"

	"github.com/VixsTy/grimoire"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

const executableName = "website"
const dockerRepository = "alexheld/website"
const outputPath = "dist"

// environmentVars holds variables used by multiple steps
var environmentVars map[string]string
var dockerImageLatest = "alexheld/website:latest"
var dockerImage = "alexheld/website:preview"
var dockerTag = "preview"
var outputDirectory = path.Join("dist", executableName)

func init() {
	dockerTag = grimoire.Tag()
	dockerImage = Sprintf("%s:%s", dockerRepository, dockerTag)
	environmentVars = map[string]string{
		"TAG":          dockerTag,
		"DOCKER_IMAGE": dockerImage,
	}
}

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	Println("## Building")
	if err := Get(); err != nil {
		return err
	}
	cmd := exec.Command("go", "build", "-o", outputDirectory, "./...")
	return cmd.Run()
}

// A formatting step that fixes alignment smaller issues and makes sure that
func Fmt() error {
	Println("## Formatting Code")
	return sh.Run("gofmt", "-w", "-s", ".")
}

// A build step that requires additional params, or platform specific steps for example
func Dockerbuild() error {
	Println("## Building Docker Image using tag 'latest'")
	err := sh.RunV("docker-compose", "-f", "docker-compose.yaml", "build")
	return err
}

// A build step that requires additional params, or platform specific steps for example
func Dockerpublish() error {
	if err := Dockerbuild(); err != nil {
		return err
	}

	Println("## Logging into DockerHub")
	err := sh.RunWithV(environmentVars, "docker", "login", "-u", "$DOCKER_USERNAME", "--password=$DOCKER_PASSWORD")
	if err != nil {
		return err
	}

	Println(Sprintf("## Pushing image to DockerHub"))
	err = sh.RunV("docker", "push", dockerRepository)
	if err != nil {
		return err
	}

	Println(Sprintf("## Tagging image as '%s'", environmentVars["TAG"]))
	err = sh.RunWithV(environmentVars, "docker", "tag", dockerImageLatest, dockerImage)
	if err != nil {
		return err
	}

	Println(Sprintf("## Pushing image to DockerHub"))
	err = sh.RunWithV(environmentVars, "docker-compose", "push")
	return err
}

// A build step that prints the current version tag
// see: getTagInternal
func Tag() error {
	tag := grimoire.Tag()
	Printf("## Resolving Tag\n")
	println(tag)
	return nil
}

// A build step that requires additional params, or platform specific steps for example
func Get() error {
	Println("Get go module dependencies...")
	cmd := exec.Command("go", "get", "-v", "-u", "-t", "./...")
	return cmd.Run()
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Build)

	installFunc := func(installDir string) error {
		Printf("Installing %s into %s", executableName, installDir)
		return os.Rename(outputDirectory, path.Join(installDir, executableName))
	}
	if gobin := os.Getenv("GOBIN"); gobin != "" {
		return installFunc(gobin)
	}
	if gopath := os.Getenv("GOPATH"); gopath != "" {
		return installFunc(path.Join(gopath, "bin"))
	}
	if goroot := os.Getenv("GOROOT"); goroot != "" {
		return installFunc(path.Join(goroot, "bin"))
	}

	return Errorf("unable to locate neither GOPATH, GOBIN nor GOROOT to resolve the installation directory")
}

func currentDir() string {
	workdir, err := os.Getwd()
	if err != nil {
		workdir = "."
	}
	return workdir
}

// Clean up after yourself
func Clean() error {
	err := grimoire.Go{}.Tidy()
	if err != nil {
		return err
	}
	Println("## Cleaning output directory")
	return os.RemoveAll(Sprintf("./%v", outputPath))
}

//  Run go vet linter
func Vet() error {
	if err := sh.Run("go", "vet", "./..."); err != nil {
		return Errorf("error running go vet: %v", err)
	}
	return nil
}
