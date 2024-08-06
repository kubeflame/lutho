//go:build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/urfave/cli/v2"
)

var platforms = []string{
	"linux/arm64",
	"linux/amd64",
	"darwin/arm64",
	"darwin/amd64",
	"windows/arm64",
	"windows/amd64",
}

func cmdRun(env []string, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, env...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println("$ ", cmd.String())
	return cmd.Run()
}

type Build struct {
	Name, AppVersion, GOOS, GOARCH string
}

func (b Build) Execute() error {
	envVars := []string{"GOOS=" + b.GOOS, "GOARCH=" + b.GOARCH}
	outputName := fmt.Sprintf("%s-%s-%s", b.Name, b.GOOS, b.GOARCH)
	ldflags := fmt.Sprintf("-ldflags=-s -w -X 'main.AppVersion=%s'", b.AppVersion)
	if b.GOOS == "windows" {
		outputName += ".exe"
	}
	err := cmdRun(envVars, "go", "build", ldflags, "-o", outputName, b.Name)
	if err != nil {
		return err
	}
	return nil
}

var versionFlagValue string
var versionFlag = &cli.StringFlag{
	Required:    true,
	Name:        "version",
	Usage:       "application version",
	Destination: &versionFlagValue,
}

var nameFlagValue string
var nameFlag = &cli.StringFlag{
	Required:    true,
	Name:        "name",
	Usage:       "application name",
	Destination: &nameFlagValue,
}

func main() {
	app := &cli.App{
		Authors: []*cli.Author{
			{
				Name:  "KubeFlame",
				Email: "https://github.com/kubeflame",
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "build",
				Usage:       "build the application",
				Description: "",
				Flags:       []cli.Flag{nameFlag, versionFlag},
				Subcommands: []*cli.Command{
					{
						Name: "version",
						Action: func(ctx *cli.Context) error {
							return nil
						},
					},
					{
						Name: "name",
						Action: func(ctx *cli.Context) error {
							return nil
						},
					},
				},
				Action: func(ctx *cli.Context) error {
					tasks := map[string]func() error{}

					err := cmdRun(os.Environ(), "go", "generate", "./...")
					if err != nil {
						return err
					}

					for key, platform := range platforms {
						splitPlatfrom := strings.Split(platforms[key], "/")
						b := Build{
							Name:       nameFlagValue,
							AppVersion: versionFlagValue,
							GOOS:       splitPlatfrom[0],
							GOARCH:     splitPlatfrom[1],
						}
						tasks[platform] = b.Execute
					}

					for taskName := range tasks {
						taskRun, ok := tasks[taskName]
						if !ok {
							log.Println("error: no such task:", taskName)
						}
						err := taskRun()
						if err != nil {
							log.Println(err)
						}
					}

					return nil
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
