package services

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/ssh"
)

const (
	doKey     = ""
	sshKey    = 0
	masterKey = ""
)

func handleDeploy(c *fiber.Ctx) error {
	return nil
}

func deploy(id string, apiKey string, region string, size string) error {
	client := godo.NewFromToken(doKey)
	ctx := context.TODO()
	req := &godo.DropletCreateRequest{
		Name:   id,
		Region: region,
		Size:   size,
		Image: godo.DropletCreateImage{
			Slug: "ubuntu-20-04-x64",
		},
		SSHKeys: []godo.DropletCreateSSHKey{},
	}
	droplet, _, err := client.Droplets.Create(ctx, req)
	if err != nil {
		return err
	}
	ip, err := droplet.PublicIPv4()
	if err != nil {
		return err
	}
	sshClient, sshSess, err := getSsh(ip)
	if err != nil {
		return err
	}
	defer sshClient.Close()
	defer sshSess.Close()
	err = install(sshSess, apiKey)
	if err != nil {
		return err
	}
	err = run(sshSess)
	if err != nil {
		return err
	}
	return nil
}

func getSsh(ip string) (*ssh.Client, *ssh.Session, error) {
	sshConfig := &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", ip+":22", sshConfig)
	if err != nil {
		return nil, nil, err
	}
	session, err := client.NewSession()
	if err != nil {
		return nil, nil, err
	}
	return client, session, nil
}

func install(sess *ssh.Session, apiKey string) error {
	cmd := fmt.Sprintf(`
		export HEISENBERG_MASTER_KEY=%s && 
		export HEISENBERG_API_KEY=%s
		export CGO_ENABLED=1`,
		masterKey, apiKey)
	err := sess.Run(cmd)
	if err != nil {
		return err
	}
	cmd = "sudo apt update && sudo apt upgrade && sudo apt install gcc -y && sudo apt install g++ -y && sudo snap install go --classic"
	err = sess.Run(cmd)
	if err != nil {
		return err
	}
	cmd = "git clone https://github.com/quantanotes/heisenberg.git"
	err = sess.Run(cmd)
	if err != nil {
		return err
	}
	return nil
}

func run(sess *ssh.Session) error {
	cmd := "cd ~/heisenberg && go run ."
	err := sess.Run(cmd)
	if err != nil {
		return err
	}
	return nil
}
