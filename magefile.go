//+build mage

package main

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/magefile/mage/sh"
)

func BuildWindowsAmd64() error {
	return build("windows", "amd64")
}

func BuildMacAmd64() error {
	return build("darwin", "amd64")
}

func BuildMacArm64() error {
	return build("darwin", "arm64")
}

func BuildLinuxAmd64() error {
	return build("linux", "amd64")
}

func BuildLinuxArm64() error {
	return build("linux", "arm64")
}

func Build() error {
	return build(runtime.GOOS, runtime.GOARCH)
}

func BuildAll() {
	parallelBuild([](func() error){
		BuildWindowsAmd64,
		BuildMacAmd64,
		BuildMacArm64,
		BuildLinuxAmd64,
		BuildLinuxArm64,
	})
}

func Check() error {
	run("go", []string{"vet", "."}, map[string]string{})
	run("go", []string{"lint", "."}, map[string]string{})
	return nil
}

func Clean() error {
	fmt.Println("Removing bin")
	return sh.Rm("bin")
}

func parallelBuild(builders [](func() error)) {
	var wg sync.WaitGroup

	for _, builder := range builders {
		wg.Add(1)
		go (func(builder (func() error), wg *sync.WaitGroup) {
			defer wg.Done()
			builder()
		})(builder, &wg)
	}
	wg.Wait()
}

func build(os, arch string) error {
	extension := ""
	if os == "windows" {
		extension = ".exe"
	}

	output, err := run("go", []string{
		"build",
		"-o", fmt.Sprintf("bin/%s-%s/serious%s", os, arch, extension),
		".",
	}, map[string]string{
		"GOOS":   os,
		"GOARCH": arch,
	})
	fmt.Print(output)
	return err
}

func run(program string, args []string, env map[string]string) (string, error) {
	// Make string representation of command
	fullArgs := append([]string{program}, args...)
	cmdStr := strings.Join(fullArgs, " ")

	// Make string representation of environment
	envStrBuf := new(bytes.Buffer)
	for key, value := range env {
		fmt.Fprintf(envStrBuf, "%s=\"%s\", ", key, value)
	}
	envStr := string(bytes.TrimRight(envStrBuf.Bytes(), ", "))

	// Show info
	fmt.Println("Running '" + cmdStr + "'" + " with env " + envStr)

	// Run
	return sh.OutputWith(env, program, args...)
}
