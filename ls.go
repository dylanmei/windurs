package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/dylanmei/windurs/windurs"
	"github.com/mitchellh/cli"
)

type lsCommand struct {
	shutdown <-chan struct{}
	ui       cli.Ui
}

func (c *lsCommand) Help() string {
	text := `
Usage: windurs ls [options] directory

  List files of a remote directory.

Options:

  -addr=localhost:5985    Host and port of the remote machine
  -user                   Name of the user to authenticate as
  -pass                   Password to authenticate with
`
	return strings.TrimSpace(text)
}

func (c *lsCommand) Synopsis() string {
	return "List files of a remote directory"
}

func (c *lsCommand) Run(args []string) int {
	var user string
	var pass string
	var addr string

	flags := flag.NewFlagSet("ls", flag.ContinueOnError)
	flags.Usage = func() { c.ui.Output(c.Help()) }
	flags.StringVar(&user, "user", "", "auth name")
	flags.StringVar(&pass, "pass", "", "auth password")
	flags.StringVar(&addr, "addr", "localhost:5985", "remote addr")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) < 1 {
		c.ui.Error("A directory is required.\n")
		c.ui.Error(c.Help())
		return 1
	}
	if len(args) > 1 {
		c.ui.Error("Too many arguments. Only a directory is required.\n")
		c.ui.Error(c.Help())
		return 1
	}

	dir := args[0]

	w, err := windurs.New(addr, user, pass)
	if err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	list, err := w.List(dir)
	if err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	tw := new(tabwriter.Writer)
	tw.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintf(tw, "\tName\tLastWriteTime\tLength\n")
	fmt.Fprintf(tw, "\t----\t-------------\t------\n")
	for _, fi := range list {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%d\n", fi.Mode, fi.Name, fi.LastWriteTime, fi.Length)
	}
	tw.Flush()

	return 0
}
