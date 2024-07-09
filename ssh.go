package main

import (
	"io"
	"os"

	"golang.org/x/crypto/ssh"
)

func uploadFileToRemote(filename string, remotePath string, username string, password string) error {
	// SSH-Verbindungsdetails
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// SSH-Verbindung herstellen
	client, err := ssh.Dial("tcp", "192.168.3.109:22", sshConfig) // IP-Adresse angepasst
	if err != nil {
		return err
	}
	defer client.Close()

	// Erstellen einer neuen SSH-Sitzung
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// Öffnen der lokalen Datei
	localFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer localFile.Close()

	// Erstellen eines Pipes zum Übertragen der Datei
	pipeReader, pipeWriter := io.Pipe()

	go func() {
		defer pipeWriter.Close()
		io.Copy(pipeWriter, localFile)
	}()

	session.Stdin = pipeReader
	if err := session.Run("/usr/bin/scp -t " + remotePath); err != nil {
		return err
	}

	return nil
}
