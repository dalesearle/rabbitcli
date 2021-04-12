package cmd

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func promptUserYN(prompt string) bool {
	var ans string

	fmt.Print(prompt + " (y/n)? ")
	fmt.Scan(&ans)
	return strings.EqualFold(ans,"y")
}



func promptUser(prompt string) string {
	var ans string

	fmt.Print(prompt)
	fmt.Scan(&ans)
	return ans
}

func promptUserSecret(prompt string) []byte {
	// Get the initial state of the terminal.
	initialTermState, e1 := terminal.GetState(syscall.Stdin)
	if e1 != nil {
		panic(e1)
	}

	// Restore it in the event of an interrupt.
	// CITATION: Konstantin Shaposhnikov - https://groups.google.com/forum/#!topic/golang-nuts/kTVAbtee9UA
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		_ = terminal.Restore(syscall.Stdin, initialTermState)
		os.Exit(1)
	}()

	// Now get the password.
	fmt.Print(prompt)
	p, err := terminal.ReadPassword(syscall.Stdin)
	fmt.Println("")
	if err != nil {
		panic(err)
	}

	// Stop looking for ^C on the channel.
	signal.Stop(c)
	return p
}
