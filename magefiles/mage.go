package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

var Default = Build

func Build(ctx context.Context) error {
	return run(ctx, makeTemp, buildWasm, buildApp)
}

func Deploy(ctx context.Context) error {
	err := Build(ctx)
	if err != nil {
		return err
	}
	err = exec.Command("scp", "temp/app.bin", "optiplex:/tmp/asteroids").Run()
	if err != nil {
		return err
	}
	return exec.Command("ssh", "optiplex", "sudo", "/tmp/asteroids", "-action", "deploy").Run()
}

func makeTemp(ctx context.Context) (err error) {
	return os.MkdirAll("temp", 0755)
}

func buildWasm(ctx context.Context) error {
	return goCmd(true, "build", "-o", "web/app.wasm", "main.go")
}

func buildApp(ctx context.Context) error {
	return goCmd(false, "build", "-o", "temp/app.bin", "web/server.go")
}

func goCmd(wasm bool, args ...string) error {
	cmd := exec.Command("go", args...)
	if wasm {
		cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	} else {
		cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	}
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		fmt.Println(string(output))
	}
	return err
}

func run(ctx context.Context, commands ...func(ctx context.Context) error) error {
	for _, command := range commands {
		if err := command(ctx); err != nil {
			return err
		}
	}
	return nil
}
