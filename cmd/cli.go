package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/lox/bkl/runner"
)

type CLI struct {
	Debug           bool              `help:"Whether to show logging"`
	DebugHTTP       bool              `help:"Whether to show http logging"`
	File            *os.File          `help:"The buildkite pipeline file to read"`
	Env             []string          `help:"Environment variables to set for the build"`
	Metadata        map[string]string `help:"Metadata to set for the build"`
	Command         string            `help:"A command to execute"`
	StepFilterRegex string            `help:"A regex to filter which steps are run"`
	Prompt          bool              `help:"Prompt before running steps"`
	DryRun          bool              `help:"Dry run, don't actually run"`
	ListenPort      int               `help:"The port to run on"`
	BootstrapScript string            `help:"Specify an alternative custom bootstrap command" type:"existingfile"`
}

func (c *CLI) Run(ctx *kong.Context) error {
	if c.Debug {
		runner.Debug = true
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	cancelCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-quit
		fmt.Printf("\n>>> Gracefully shutting down...\n")
		cancel()
	}()

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	commit, err := gitCommit()
	if err != nil {
		log.Printf("Error getting git commit: %v", err)
		commit = "no_commit_found"
	}

	branch, err := gitBranch()
	if err != nil {
		log.Printf("Error getting git branch: %v", err)
		branch = "master"
	}

	command := c.Command
	if c.File != nil {
		command = fmt.Sprintf("buildkite-agent pipeline upload %q", c.File.Name())
	}

	stepFilterReg, err := regexp.Compile(c.StepFilterRegex)
	if err != nil {
		return err
	}

	return runner.Run(cancelCtx, runner.Params{
		Debug:           c.Debug,
		DebugHTTP:       c.DebugHTTP,
		Env:             c.Env,
		Metadata:        c.Metadata,
		DryRun:          c.DryRun,
		Command:         command,
		Dir:             wd,
		Prompt:          c.Prompt,
		StepFilter:      stepFilterReg,
		ListenPort:      c.ListenPort,
		BootstrapScript: c.BootstrapScript,
		JobTemplate: runner.Job{
			Commit:           commit,
			Branch:           branch,
			Repository:       wd,
			OrganizationSlug: "local",
			PipelineSlug:     filepath.Base(wd),
		},
	})
}

func gitBranch() (string, error) {
	out, err := exec.Command(`git`, `rev-parse`, `--abbrev-ref`, `HEAD`).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func gitCommit() (string, error) {
	out, err := exec.Command(`git`, `rev-parse`, `HEAD`).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
