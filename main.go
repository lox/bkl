package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/lox/bkl/cmd"
)

var (
	version string // set by goreleaser
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Shutting down")
		signal.Stop(c)
		cancel()
	}()

	if err := run(ctx); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cli := cmd.CLI{}

	k := kong.Parse(&cli,
		kong.Name("bkl"),
		kong.Description("Run your Buildkite pipelines locally (no buildkite.com interaction)"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": version,
		})

	k.BindTo(ctx, (*context.Context)(nil))
	return k.Run(&cli)
}
