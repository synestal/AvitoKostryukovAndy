package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	if _, err := exec.LookPath("fieldalignment"); err != nil {
		if err := installFieldalignment(); err != nil {
			log.Fatal(err)
		}
	}
	if _, err := exec.LookPath("golangci-lint"); err != nil {
		if err := installGolangciLint(); err != nil {
			log.Fatal(err)
		}
	}
	if err := os.Chdir("./.."); err != nil {
		log.Fatalf("err: %v", err)
	}
	if err := runGolangciLint(); err != nil {
		log.Println("State", err)
	}
	if err := os.Chdir("./cmd/avito-test-trainee"); err != nil {
		log.Fatalf("err:", err)
	}
	if err := runFieldalignment(); err != nil {
		log.Printf("err:", err)
	}
}

func installGolangciLint() error {
	cmd := exec.Command("go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runGolangciLint() error {
	cmd := exec.Command("golangci-lint", "run")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installFieldalignment() error {
	cmd := exec.Command("go", "install", "golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runFieldalignment() error {
	cmd := exec.Command("fieldalignment", "-fix", "./")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
