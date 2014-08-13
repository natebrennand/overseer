package main

import (
	"github.com/natebrennand/overseer/output"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	Delay = time.Second
)

var watchFiles []string
var filePatterns []regexp.Regexp
var commandArgs []string

func matchFile(path string, info os.FileInfo, err error) error {
	if err != nil || path[0] == '.' { //|| info.Name()[0] != '.'{
		return nil
	}

	if !info.IsDir() {
		for _, pattern := range filePatterns {
			if pattern.MatchString(path) {
				watchFiles = append(watchFiles, path)
				return nil
			}
		}
	}
	return nil
}

func findFiles(patterns []string) {
	filePatterns = []regexp.Regexp{}

	for _, pat := range patterns {
		if strings.Contains(pat, "**/*") {
			// let **/* access all directories
			strings.Replace(pat, "**/*", "*", 50)
		} else if strings.Contains(pat, "*") {
			// prevent wildcard from accessing other directories
			strings.Replace(pat, "*", "[^/]*", 50)
		}
		filePatterns = append(filePatterns, *(regexp.MustCompile("^" + pat)))
	}
	watchFiles = []string{}
	filepath.Walk("./", matchFile)
}

// parse CLI args for files to watch and the command to execute
func parseComands() {
	if len(os.Args) < 4 {
		output.Usuage()
	}

	commandIndex := -1
	for i, w := range os.Args {
		if w == "-c" {
			commandIndex = i
			break
		}
	}
	if commandIndex < 0 {
		output.Usuage()
	}
	commandArgs = os.Args[commandIndex+1:]
	findFiles(os.Args[1:commandIndex])
}

// Initializes the last time every file was modified
func initFilesModTimes() map[string]time.Time {
	fileModTimes := make(map[string]time.Time)

	for _, f := range watchFiles {
		fInfo, err := os.Stat(f)
		if err != nil {
			output.FatalError(err)
		}
		fileModTimes[f] = fInfo.ModTime()
	}

	return fileModTimes
}

// Determines if any files were modified
// updates their lastModified time
func filesModified(fileModTimes map[string]time.Time) bool {
	returnVal := false
	for f := range fileModTimes {
		fInfo, err := os.Stat(f)
		if err != nil {
			output.FatalError(err)
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
	var commandOutput bytes.Buffer
	var c *exec.Cmd
	if len(commandArgs) > 1 {
		c = exec.Command(commandArgs[0], commandArgs[1:]...)
	} else {
		c = exec.Command(commandArgs[0])
	}
	c.Stdout = &commandOutput
	c.Stderr = &commandOutput

	err := c.Run()
	if err != nil {
		output.PrintError(commandOutput, err)
	} else if commandOutput.String() != "" {
		output.PrintSuccess(commandOutput)
	} else {
		output.NoError()
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
