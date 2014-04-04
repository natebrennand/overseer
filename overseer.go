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

var watchFiles []string
var commandArgs []string

func parseComands() {
	commandArgs = []string{}
	watchFiles = []string{}

	onCommands := false
	for _, w := range(os.Args) {
		if w == "-c" {
			onCommands = true
			continue
		}

		if onCommands {
			commandArgs = append(commandArgs, w)
		} else {
			watchFiles = append(watchFiles, w)
		}
	}
}

func initFilesModTimes() map[string]time.Time {
	fileModTimes := make(map[string]time.Time)

	for _, f := range(watchFiles) {
		fInfo, err := os.Stat(f)
		if err != nil {
			log.Fatal(err.Error())
		}
		fileModTimes[f] = fInfo.ModTime()
	}

	return fileModTimes
}

func filesModified (fileModTimes map[string]time.Time) bool {
	returnVal := false
	for f := range(fileModTimes) {
		fInfo, err := os.Stat(f)
		if err != nil {
			log.Fatal(err.Error())
		}
		if fileModTimes[f] != fInfo.ModTime() {
			fileModTimes[f] = fInfo.ModTime()
			returnVal = true
		}
	}
	return returnVal
}

// run whatever command was passed in when started
func runCommand() {
	var output bytes.Buffer
	var c *exec.Cmd
	if len(commandArgs) > 1 {
		c = exec.Command(commandArgs[0], commandArgs[1:]...)
	} else {
		c = exec.Command(commandArgs[0])
	}
	c.Stdout = &output
	c.Stderr = &output

	err := c.Run()
	if err != nil {
		fmt.Println(output.String())
		log.Fatal(err.Error())
	}

	if output.String() != "" {
		fmt.Println(output.String())
	}
}

func main() {
	parseComands()
	fileModTimes := initFilesModTimes()

	for {
		time.Sleep(Delay)
		if filesModified(fileModTimes) {
			runCommand()
		}
	}
}
