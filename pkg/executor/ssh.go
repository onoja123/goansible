package executor

import (
	"fmt"
	"io/ioutil"
	"time"

	"goansible/pkg/inventory"

	"golang.org/x/crypto/ssh"
)

type SSHConnection struct {
	client *ssh.Client
	host   *inventory.Host
}

func NewSSHConnection(host *inventory.Host) (*SSHConnection, error) {
	var authMethods []ssh.AuthMethod

	if host.KeyFile != "" {
		key, err := ioutil.ReadFile(host.KeyFile)
		if err != nil {
			return nil, err
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, err
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if host.Password != "" {
		authMethods = append(authMethods, ssh.Password(host.Password))
	}

	config := &ssh.ClientConfig{
		User:            host.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", host.Address, host.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	return &SSHConnection{
		client: client,
		host:   host,
	}, nil
}

func (c *SSHConnection) Execute(command string) (string, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	return string(output), err
}

func (c *SSHConnection) ExecuteWithSudo(command, becomeUser string) (string, error) {
	if becomeUser == "" {
		becomeUser = "root"
	}
	sudoCmd := fmt.Sprintf("sudo -u %s %s", becomeUser, command)
	return c.Execute(sudoCmd)
}

func (c *SSHConnection) Close() error {
	return c.client.Close()
}

func Ping(host *inventory.Host) error {
	conn, err := NewSSHConnection(host)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Execute("echo pong")
	return err
}

func GatherFacts(host *inventory.Host) (map[string]interface{}, error) {
	conn, err := NewSSHConnection(host)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	facts := make(map[string]interface{})

	// Gather basic facts
	if hostname, err := conn.Execute("hostname"); err == nil {
		facts["hostname"] = hostname
	}

	if os, err := conn.Execute("uname -s"); err == nil {
		facts["os"] = os
	}

	return facts, nil
}
