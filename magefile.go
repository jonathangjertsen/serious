//+build mage

package main

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
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

func Ci() {
	mg.Deps(Check)
	mg.Deps(Test)
	mg.Deps(BuildAll)
}

func Check() error {
	output, err := run("go", []string{"vet", "./..."}, map[string]string{})
	if err != nil {
		return err
	}
	if output != "" {
		return fmt.Errorf("go vet says something:\n%s", output)
	}

	output, err = run("go", []string{"fmt", "./..."}, map[string]string{})
	if err != nil {
		return err
	}
	if output != "" {
		return fmt.Errorf("go fmt says something:\n%s", output)
	}
	return nil
}

func Test() error {
	output, err := run("go", []string{"test", "-v", "./..."}, map[string]string{})
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "[no test files]") {
			continue
		}
		if strings.HasPrefix(line, "ok") {
			color.HiGreen(line)
		} else if strings.Contains(line, "FAIL") {
			color.HiRed(line)
		} else {
			fmt.Println(line)
		}
	}
	return err
}

func Clean() error {
	fmt.Println("Removing bin")
	return sh.Rm("bin")
}

func parallelBuild(builders [](func() error)) {
	var wg sync.WaitGroup

	for _, builder := range builders {
		wg.Add(1)
		go (func(builder func() error, wg *sync.WaitGroup) {
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
