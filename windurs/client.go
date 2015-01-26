package windurs

import (
	"io"

	"github.com/masterzen/winrm/winrm"
)

type Client struct {
	inner *winrm.Client
}

func New(addr, user, pass string) (*Client, error) {
	endpoint, err := parseEndpoint(addr)
	if err != nil {
		return nil, err
	}
	inner := winrm.NewClient(endpoint, user, pass)
	return &Client{inner}, nil
}

func (c *Client) Info() (*Info, error) {
	return fetchInfo(c.inner)
}

func (c *Client) List(remotePath string) ([]FileItem, error) {
	return fetchList(c.inner, winPath(remotePath))
}

func (c *Client) Cmd(stdout, stderr io.Writer, impersonate bool, arguments ...string) (int, error) {
	if !impersonate {
		return runCmd(c.inner, stdout, stderr, arguments...)
	} else {
		return runElevatedCmd(c.inner, stdout, stderr, arguments...)
	}
}

