package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"bufio"
	"time"
	"strings"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"github.com/mxseba/go-term"
)

func Terminal(port int) {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("login as: ")
	username, _ := reader.ReadString('\n')
	fmt.Print("password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		log.Println(err)
	}
	password := strings.TrimSpace(string(bytePassword))
	auth := []ssh.AuthMethod{ssh.Password(password)}
	config := &ssh.ClientConfig{
		User: strings.TrimSpace(username),
		Auth: auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: 30 * time.Second,
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port), config)
	
	
	if err != nil {
		fmt.Println(err) 
	} 
	session, err := client.NewSession()
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close() 
			
	session.Stdin, session.Stdout, session.Stderr = term.StdStreams()
			
	fdIn := int(syscall.Stdin)
	_, err = terminal.MakeRaw(fdIn)
	if err != nil {
		fmt.Println(err)
	}
	
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.ECHOCTL:       1,	  // enable echoing ctl
		ssh.IGNCR: 		   0, 	  // Ignore CR on input.
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	
	fdOut := int(os.Stdout.Fd())
	
	termWidth, termHeight, err := terminal.GetSize(fdOut)
	if err != nil {
		fmt.Println(err)
	}
	err = session.RequestPty("xterm-256color", termHeight, termWidth, modes)
	
	if err != nil {
		fmt.Println(err)
	}
	err = session.Shell()
	if err != nil {
		fmt.Println(err)
	}
	
	
	go resizeEvent(session, termWidth, termHeight)
	session.Wait()

}



func resizeEvent(session *ssh.Session, termWidth, termHeight int){
	for{
		time.Sleep(100 * time.Millisecond)
		fdOut := int(os.Stdout.Fd())
		newTermWidth, newTermHeight, err := terminal.GetSize(fdOut)
		if err != nil {
			fmt.Println(err)
		}
		if termWidth != newTermWidth || termHeight != newTermHeight {
			err = session.WindowChange(newTermHeight, newTermWidth)
			if err != nil {
				fmt.Println(err)
			}
			termWidth = newTermWidth
			termHeight = newTermHeight
		}
	}
}

