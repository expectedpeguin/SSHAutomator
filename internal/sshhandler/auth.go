package sshhandler

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

func GetAuthMethod(password, keyFile string) (ssh.AuthMethod, error) {
	if password != "" {
		return ssh.Password(password), nil
	} else if keyFile != "" {
		key, err := ioutil.ReadFile(keyFile)
		if err != nil {
			return nil, err
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, err
		}
		return ssh.PublicKeys(signer), nil
	}
	return nil, fmt.Errorf("either -password or -keyfile flag must be provided for authentication")
}
