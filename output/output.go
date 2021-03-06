package output

import (
	"bytes"
	"fmt"
	"log"
)

func PrintError(b bytes.Buffer, err error) {
	redBackground()
	log.Printf("ERROR:\033[0m %s", err.Error())
	fmt.Print(b.String())
}

func PrintSuccess(b bytes.Buffer) {
	NoError()
	fmt.Print(b.String())
}

func NoError() {
	greenColor()
	log.Printf("No error")
	endColor()
}

func Usuage() {
	fmt.Println("Usuage:\t./overseer <file> ... -c command ...")
	redBackground()
	log.Fatal("Incorrect usuage\033[0m")
}

func FatalError(err error) {
	redBackground()
	log.Fatalf("\033[0m%s", err.Error())
}

func greenColor() {
	fmt.Printf("\033[32m")
}

func redBackground() {
	fmt.Printf("\033[1;41m")
}

func endColor() {
	fmt.Printf("\033[0m")
}
