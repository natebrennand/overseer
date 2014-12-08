package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/howeyc/fsnotify"
	"github.com/natebrennand/overseer/output"
)

var filePatterns []regexp.Regexp

func matchFile(path string, info os.FileInfo, err error) error {
	if err != nil || path[0] == '.' {
		return nil
	}

	if !info.IsDir() {
		for _, pattern := range filePatterns {
			if pattern.MatchString(path) {
				// watchFiles = append(watchFiles, path)
				return nil
			}
		}
	}
	return nil
}

func findFiles(patterns []string) []string {
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
	watchFiles := []string{}
	filepath.Walk("./", filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
		if err != nil || path[0] == '.' {
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
	}))

	return watchFiles
}

// run whatever command was passed in when started
func runCommand(command []string) {
	var commandOutput bytes.Buffer
	var c *exec.Cmd
	if len(command) > 1 {
		c = exec.Command(command[0], command[1:]...)
	} else {
		c = exec.Command(command[0])
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

// parse CLI args for files to watch and the command to execute
func parseComands() ([]string, []string) {
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
	commandArgs := os.Args[commandIndex+1:]
	return commandArgs, findFiles(os.Args[1:commandIndex])
}

func main() {
	done := make(chan bool)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	// parser CLI args and get files list
	command, files := parseComands()

	// add all files to watcher
	for _, f := range files {
		log.Println("Watching:", f)
		// if err = watcher.WatchFlags(f, fsnotify.FSN_RENAME|fsnotify.FSN_MODIFY); err != nil {
		if err = watcher.WatchFlags(f, fsnotify.FSN_ALL); err != nil {
			log.Fatal(err)
		}
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				log.Println(ev)
				runCommand(command)
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	<-done
	watcher.Close()
}
