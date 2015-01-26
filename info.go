package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/dylanmei/windurs/windurs"
	"github.com/mitchellh/cli"
)

type infoCommand struct {
	shutdown <-chan struct{}
	ui       cli.Ui
}

func (c *infoCommand) Help() string {
	text := `
Usage: winrmfs info [options]

  Show status and info about the remote.

Options:

  -addr=localhost:5985    Host and port of the remote machine
  -user=""                Name of the user to authenticate as
  -pass=""                Password to authenticate with
`
	return strings.TrimSpace(text)
}

func (c *infoCommand) Synopsis() string {
	return "Show status and info about the remote"
}

func (c *infoCommand) Run(args []string) int {
	var user string
	var pass string
	var addr string

	flags := flag.NewFlagSet("info", flag.ContinueOnError)
	flags.Usage = func() { c.ui.Output(c.Help()) }
	flags.StringVar(&user, "user", "", "auth name")
	flags.StringVar(&pass, "pass", "", "auth password")
	flags.StringVar(&addr, "addr", "localhost:5985", "remote addr")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	w, err := windurs.New(addr, user, pass)
	if err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	info, err := w.Info()
	if err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	c.ui.Output("Client")
	c.ui.Output(fmt.Sprintf("    Addr: %s", addr))
	c.ui.Output(fmt.Sprintf("    Auth: %s", "Basic"))
	c.ui.Output(fmt.Sprintf("    User: %s", user))
	c.ui.Output("WinRM Config")
	c.ui.Output(fmt.Sprintf("    %s: %d", "MaxEnvelopeSizeKB", info.WinRM.MaxEnvelopeSizeKB))
	c.ui.Output(fmt.Sprintf("    %s: %d", "MaxTimeoutMS", info.WinRM.MaxTimeoutMS))
	c.ui.Output(fmt.Sprintf("    %s: %d", "Service/MaxConcurrentOperations", info.WinRM.Service.MaxConcurrentOperations))
	c.ui.Output(fmt.Sprintf("    %s: %d", "Service/MaxConcurrentOperationsPerUser", info.WinRM.Service.MaxConcurrentOperationsPerUser))
	c.ui.Output(fmt.Sprintf("    %s: %d", "Service/MaxConnections", info.WinRM.Service.MaxConnections))
	c.ui.Output(fmt.Sprintf("    %s: %d", "Winrs/MaxConcurrentUsers", info.WinRM.Winrs.MaxConcurrentUsers))
	c.ui.Output(fmt.Sprintf("    %s: %d", "Winrs/MaxProcessesPerShell", info.WinRM.Winrs.MaxProcessesPerShell))
	c.ui.Output(fmt.Sprintf("    %s: %d", "Winrs/MaxMemoryPerShellMB", info.WinRM.Winrs.MaxMemoryPerShellMB))
	c.ui.Output(fmt.Sprintf("    %s: %d", "Winrs/MaxShellsPerUser", info.WinRM.Winrs.MaxShellsPerUser))

	return 0
}
