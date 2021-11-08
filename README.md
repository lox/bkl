# bkl - Run Buildkite Pipelines locally

Run buildkite pipelines locally for fast testing or if buildkite.com is down.

## Installing

Either download a release binary or if you have golang:

```bash
go install github.com/lox/bkl@latest
```

## Running

```bash
$ bkl

>>> Starting local agent 🤖
>>> Starting build 👟
>>> Executing initial command: buildkite-agent pipeline upload "/Users/lachlan/Projects/lox/bkl/examples/hello-world.pipeline.yml"
>>> Executing command step

~~~ Preparing plugins
$ cd /Users/lachlan/Projects/lox/bkl

~~~ Running commands
$ trap 'kill -- $$' INT TERM QUIT; echo hello world!
hello world!

>>> Command succeeded in 3.047168424s
>>> Build finished in 3.051489877s
```

## Credit

Extracted from https://github.com/buildkite/cli.

P.S. @toolmantim, you were right. This always should have been a standalone tool. Soz.
