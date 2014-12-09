package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/howeyc/fsnotify"
	"github.com/natebrennand/overseer/output"
)

// findFiles makes a list of directories to watch and a regex to match filenames with
func findFiles(patterns []string) ([]string, regexp.Regexp) {
	filePatterns := []string{}

	for _, pat := range patterns {
		if strings.Contains(pat, "**/*") { // let **/* access all directories
			strings.Replace(pat, "**/*", "*", 50)
		} else if strings.Contains(pat, "*") { // prevent wildcard from accessing other directories
			strings.Replace(pat, "*", "[^/]*", 50)
		}
		filePatterns = append(filePatterns, "^"+pat)
	}

	fileregex := regexp.MustCompile(strings.Join(filePatterns, "|"))
	set := make(map[string]bool)
	watchDirs := []string{}

	filepath.Walk("./", filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
		if err != nil || path[0] == '.' {
			return nil
		}
		dir := filepath.Dir(path)

		if fileregex.MatchString(path) {
			_, exists := set[dir]
			if exists == false {
				set[dir] = true
				watchDirs = append(watchDirs, dir)
			}
		}
		return nil
	}))

	return watchDirs, *fileregex
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
func parseComands() ([]string, []string, regexp.Regexp) {
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
	watchDirs, pattern := findFiles(os.Args[1:commandIndex])
	return commandArgs, watchDirs, pattern
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// parser CLI args and get files list
	command, files, pattern := parseComands()

	// add all directories to watcher
	for _, f := range files {
		if err = watcher.WatchFlags(f, fsnotify.FSN_MODIFY); err != nil {
			log.Fatal(err)
		}
	}

	var lastModifyTime time.Time = time.Now()
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if pattern.MatchString(ev.Name) {
					// limit calling the command
					if time.Now().Sub(lastModifyTime) > 2*time.Second {
						go runCommand(command)
						lastModifyTime = time.Now()
					}
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	done := make(chan bool)
	<-done
}
