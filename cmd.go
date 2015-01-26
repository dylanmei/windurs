package main

import (
	"flag"
	"os"
	"strings"

	"github.com/dylanmei/windurs/windurs"
	"github.com/mitchellh/cli"
)

type cmdCommand struct {
	shutdown <-chan struct{}
	ui       cli.Ui
}

func (c *cmdCommand) Help() string {
	text := `
Usage: windurs cmd [options] command

  Execute a command against cmd.exe

Options:

  -addr=localhost:5985    Host and port of the remote machine
  -user=""                Name of the user to authenticate as
  -pass=""                Password to authenticate with
  -impersonate            If true, impersonate as if logged-in (slow)
`
	return strings.TrimSpace(text)
}

func (c *cmdCommand) Synopsis() string {
	return "Execute a command against cmd.exe"
}

func (c *cmdCommand) Run(args []string) int {
	var user string
	var pass string
	var addr string
	var impersonate bool

	flags := flag.NewFlagSet("cmd", flag.ContinueOnError)
	flags.Usage = func() { c.ui.Output(c.Help()) }
	flags.StringVar(&user, "user", "", "auth name")
	flags.StringVar(&pass, "pass", "", "auth password")
	flags.StringVar(&addr, "addr", "localhost:5985", "remote addr")
	flags.BoolVar(&impersonate, "impersonate", false, "")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) < 1 {
		c.ui.Error("A command is required.\n")
		c.ui.Error(c.Help())
		return 1
	}

	w, err := windurs.New(addr, user, pass)
	if err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	exitCode, err := w.Cmd(os.Stdout, os.Stderr, impersonate, args...)
	if err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	return exitCode
}
