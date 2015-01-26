package windurs

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/masterzen/winrm/winrm"
)

type Info struct {
	WinRM      WinrmConfig
	PowerShell PsSettings
}

type WinrmConfig struct {
	MaxEnvelopeSizeKB int `xml:"MaxEnvelopeSizekb"`
	MaxTimeoutMS      int `xml:"MaxTimeoutms"`
	Service           WinrmServiceConfig
	Winrs             WinrmWinrsConfig
}

type WinrmServiceConfig struct {
	MaxConnections                 int
	MaxConcurrentOperations        int
	MaxConcurrentOperationsPerUser int
}

type WinrmWinrsConfig struct {
	MaxMemoryPerShellMB  int
	MaxShellsPerUser     int
	MaxConcurrentUsers   int
	MaxProcessesPerShell int
}

type PsSettings struct {
	Version         string
	ExecutionPolicy string
}

func fetchInfo(client *winrm.Client, user, pass string) (*Info, error) {
	var err error
	info := &Info{
		WinRM:      WinrmConfig{},
		PowerShell: PsSettings{},
	}

	err = runPsVersion(client, &info.PowerShell)
	if err != nil {
		return info, err
	}

	err = runPsExecutionPolicy(client, &info.PowerShell)
	if err != nil {
		return info, err
	}

	err = runWinrmConfig(client, user, pass, &info.WinRM)
	if err != nil {
		return info, err
	}

	return info, nil
}

func runWinrmConfig(client *winrm.Client, user, pass string, config *WinrmConfig) error {
	command := "winrm get winrm/config -format:xml"

	stdout := bytes.NewBuffer(make([]byte, 0))
	stderr := bytes.NewBuffer(make([]byte, 0))
	_, err := runElevatedCmd(client, user, pass, stdout, stderr, command)

	if os.Getenv("WINDURS_DEBUG") != "" && stderr.Len() > 0 {
		log.Printf("STDERR returned: %s\n", stderr.String())
	}

	if err == nil && stdout.Len() > 0 {
		err = xml.Unmarshal(stdout.Bytes(), config)
	}

	return err
}

func runPsVersion(client *winrm.Client, settings *PsSettings) error {
	script := "$PSVersionTable.PSVersion | ConvertTo-Xml -NoTypeInformation -As String"
	stdout, stderr, err := client.RunWithString("powershell -Command \""+script+"\"", "")

	if err != nil {
		return errors.New(fmt.Sprintf("Couldn't execute script %s: %v", script, err))
	}

	if stderr != "" {
		if os.Getenv("WWNDURS_DEBUG") != "" {
			log.Printf("STDERR returned: %s\n", stderr)
		}
	}

	if stdout != "" {
		doc := pslist{}
		err := xml.Unmarshal([]byte(stdout), &doc)
		if err != nil {
			return errors.New(fmt.Sprintf("Couldn't parse results: %v", err))
		}

		settings.Version = doc.Objects[0].Value
	}

	return nil
}

func runPsExecutionPolicy(client *winrm.Client, settings *PsSettings) error {
	script := "Get-ExecutionPolicy | Select-Object"
	stdout, stderr, err := client.RunWithString("powershell -Command \""+script+"\"", "")

	if err != nil {
		return errors.New(fmt.Sprintf("Couldn't execute script %s: %v", script, err))
	}

	if stderr != "" {
		if os.Getenv("WINRMCP_DEBUG") != "" {
			log.Printf("STDERR returned: %s\n", stderr)
		}
	}

	if stdout != "" {
		settings.ExecutionPolicy = strings.Trim(stdout, "\n")
	}

	return nil
}
