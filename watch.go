package main

import (
	"time"
	"os"
	"os/exec"
	"fmt"
	"log"
	"bytes"
)

const (
	Delay = time.Second
)

// run whatever command was passed in
func runCommand() {
	c := exec.Command(os.Args[2], os.Args[3:]...)
	var stdout bytes.Buffer
	c.Stdout = &stdout

	err := c.Run()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(stdout.String())
}

func main() {
	f := os.Args[1]
	lastModified := time.Now()

	for {
		fmt.Println("echo")
		time.Sleep(Delay)
		fi, err := os.Stat(f)
		if err != nil {
			log.Fatal(err.Error())
		}

		if lastModified != fi.ModTime() {
			lastModified = fi.ModTime()
			runCommand()
		}
	}
}
