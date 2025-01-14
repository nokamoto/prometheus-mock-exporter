//go:build mage
// +build mage

package main

import (
	"fmt"
	"os/exec"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

var Default = All

// All runs all the tasks
func All() error {
	mg.SerialDeps(Tidy, Format)
	return nil
}

// Tidy runs go mod tidy
func Tidy() error {
	fmt.Println("Tidying...")
	return exec.Command("go", "mod", "tidy").Run()
}

// Format formats all the files in the project
func Format() error {
	mg.SerialDeps(Goimports, Gofumpt)
	return nil
}

// InstallGofumpt installs the gofumpt tool
func InstallGofumpt() error {
	fmt.Println("Installing gofumpt...")
	return exec.Command("go", "install", "mvdan.cc/gofumpt@latest").Run()
}

// Gofumpt runs the gofumpt tool
func Gofumpt() error {
	mg.Deps(InstallGofumpt)
	fmt.Println("Running gofumpt...")
	return exec.Command("gofumpt", "-w", ".").Run()
}

// InstallGoimports installs the goimports tool
func InstallGoimports() error {
	fmt.Println("Installing goimports...")
	return exec.Command("go", "install", "golang.org/x/tools/cmd/goimports@latest").Run()
}

// Goimports runs the goimports tool
func Goimports() error {
	mg.Deps(InstallGoimports)
	fmt.Println("Running goimports...")
	return exec.Command("goimports", "-w", ".").Run()
}
