package windurs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/masterzen/winrm/winrm"
	"github.com/mitchellh/packer/common/uuid"
)

func runCmd(client *winrm.Client, stdout, stderr io.Writer, arguments ...string) (int, error) {
	shell, err := client.CreateShell()
	if err != nil {
		return -1, err
	}

	defer shell.Close()
	args := []string{"/c"}
	args = append(args, arguments...)
	cmd, err := shell.Execute("cmd.exe", args...)
	if err != nil {
		return -1, err
	}

	go io.Copy(stdout, cmd.Stdout)
	go io.Copy(stderr, cmd.Stderr)

	cmd.Wait()
	return cmd.ExitCode(), nil
}

func runElevatedCmd(client *winrm.Client, stdout, stderr io.Writer, arguments ...string) (int, error) {
	// generate command
	command := strings.Join(arguments, " ")
	var buffer bytes.Buffer
	err := elevatedTemplate.Execute(&buffer, elevatedOptions{
		User:            "packer",
		Password:        "packer",
		TaskDescription: "Command: " + command,
		TaskName:        fmt.Sprintf("windurs-%s", uuid.TimeOrderedUUID()),
		EncodedCommand:  psencode([]byte(command + "; exit $LASTEXITCODE")),
	})

	if err != nil {
		return -1, errors.New(fmt.Sprintf("Couldn't compile elevated command: %v", err))
	}

	shell, err := client.CreateShell()
	if err != nil {
		return -1, err
	}
	defer shell.Close()

	cmd, err := shell.Execute("powershell", "-EncodedCommand", psencode(buffer.Bytes()))
	if err != nil {
		return -1, errors.New(fmt.Sprintf("Couldn't execute elevated command: %v", err))
	}

	go io.Copy(stdout, cmd.Stdout)
	go io.Copy(stderr, cmd.Stderr)
	cmd.Wait()

	return cmd.ExitCode(), nil
}
