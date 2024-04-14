package main

import (
	"SSHAutomator/internal/sshhandler"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
)

func main() {
	hostPtr := flag.String("host", "", "SSH host")
	portPtr := flag.Int("port", 22, "SSH port")
	usernamePtr := flag.String("username", "", "SSH username")
	passwordPtr := flag.String("password", "", "SSH password")
	keyFilePtr := flag.String("keyfile", "", "Path to private key file (optional)")
	scriptFilePtr := flag.String("script", "", "Path of the script you would like to run")
	serversFilePtr := flag.String("servers", "", "File containing server details")
	flag.Parse()

	commands, err := sshhandler.ReadScriptFile(*scriptFilePtr)

	if err != nil {
		fmt.Println("Error reading specified file:", err)
		return
	}

	var servers []sshhandler.ServerDetails
	if *serversFilePtr != "" {
		var err error
		servers, err = sshhandler.ReadServersFile(*serversFilePtr)
		if err != nil {
			fmt.Println("Error reading servers file:", err)
			return
		}
	} else {
		if *hostPtr == "" || *usernamePtr == "" {
			fmt.Println("Usage: shautomator -script script.shautomator -host HOST -port PORT -username USERNAME [-password PASSWORD | -keyfile KEYFILE]")
			fmt.Println("Usage: shautomator -script script.shautomator -servers serverlist.txt")
			return
		}
		servers = append(servers, sshhandler.ServerDetails{
			Host:     *hostPtr,
			Username: *usernamePtr,
			Password: *passwordPtr,
			KeyFile:  *keyFilePtr,
		})
		err = executeSSHCommands(*hostPtr, *portPtr, *usernamePtr, *passwordPtr, *keyFilePtr, commands)
		if err != nil {
			fmt.Println("Command execution failed:", err)
		}
	}

	commands, err = sshhandler.ReadScriptFile(*scriptFilePtr)
	if err != nil {
		fmt.Println("Error reading specified file:", err)
		return
	}

	for _, server := range servers {
		go func(server sshhandler.ServerDetails) {
			err := executeSSHCommands(server.Host, *portPtr, server.Username, server.Password, server.KeyFile, commands)
			if err != nil {
				fmt.Printf("Command execution failed on %s: %v\n", server.Host, err)
			}
		}(server)
	}

	fmt.Println("Commands execution initiated on specified servers.")
	fmt.Println("Waiting for commands to complete...")
	fmt.Scanln()
}

func executeSSHCommands(host string, port int, username, password, keyFile string, commands []string) error {
	config := &ssh.ClientConfig{
		User:            username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	auth, err := sshhandler.GetAuthMethod(password, keyFile)
	if err != nil {
		return err
	}
	config.Auth = []ssh.AuthMethod{auth}

	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		return err
	}
	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}

	if err := session.Shell(); err != nil {
		return err
	}

	for _, cmd := range commands {
		if _, err := stdin.Write([]byte(cmd + "\n")); err != nil {
			return err
		}
	}

	return session.Wait()
}
