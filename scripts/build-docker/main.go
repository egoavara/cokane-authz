package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// pwd
	tag := tag()
	build(tag)
	push(tag)
}

func tag() string {
	branch, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		panic(err)
	}
	branchNorm := strings.ReplaceAll(strings.ReplaceAll(string(branch), "/", "-"), "\n", "")
	return fmt.Sprintf("ghcr.io/egoavara/cokane-authz:%s", branchNorm)
}

func build(tag string) {
	subprocess := exec.Command("docker", "buildx", "build",
		"-f", "./scripts/build-docker/Dockerfile",
		"-t", tag,
		".",
	)
	subprocess.Stdout = os.Stdout
	subprocess.Stderr = os.Stderr
	subprocess.Stdin = os.Stdin
	err := subprocess.Run()
	if err != nil {
		panic(err)
	}
}

func push(tag string) {
	subprocess := exec.Command("docker", "push", tag)
	subprocess.Stdout = os.Stdout
	subprocess.Stderr = os.Stderr
	subprocess.Stdin = os.Stdin
	err := subprocess.Run()
	if err != nil {
		panic(err)
	}
}
