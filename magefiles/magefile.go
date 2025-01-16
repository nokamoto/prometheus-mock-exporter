//go:build mage
// +build mage

package main

import (
	"fmt"
	"os/exec"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = All

// All runs all the tasks
func All() error {
	mg.SerialDeps(Tidy, Yamlfmt, BufGenerate, Format, GoTest)
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

// InstallBuf installs the buf tool
func InstallBuf() error {
	fmt.Println("Installing buf...")
	return exec.Command("go", "install", "github.com/bufbuild/buf/cmd/buf@latest").Run()
}

// BufLint runs the buf format tool
func BufFormat() error {
	mg.Deps(InstallBuf)
	fmt.Println("Running buf format...")
	return exec.Command("buf", "format", "-w", "proto").Run()
}

// BufGenerate runs the buf generate tool
func BufGenerate() error {
	mg.Deps(BufFormat)
	fmt.Println("Running buf generate...")
	return exec.Command("buf", "generate", "--clean", "--template", "build/buf.gen.yaml").Run()
}

// InstallYamlfmt installs the yamlfmt tool
func InstallYamlfmt() error {
	fmt.Println("Installing yamlfmt...")
	return exec.Command("go", "install", "github.com/google/yamlfmt/cmd/yamlfmt@latest").Run()
}

// Yamlfmt runs the yamlfmt tool
func Yamlfmt() error {
	mg.Deps(InstallYamlfmt)
	fmt.Println("Running yamlfmt...")
	return exec.Command("yamlfmt", ".").Run()
}

// GoTest runs go test
func GoTest() error {
	fmt.Println("Running go test...")
	return sh.RunV("go", "test", "./...")
}
