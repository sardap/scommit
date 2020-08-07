// This is made for me don't bitch
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	aurora "github.com/logrusorgru/aurora"
	"github.com/sardap/scommit/spellchecker"
)

func processArgs() (string, error) {
	if len(os.Args) < 2 {
		return "", fmt.Errorf("missing arg message")
	}

	var builder strings.Builder
	for _, arg := range os.Args[1:] {
		fmt.Fprintf(&builder, "%s ", arg)
	}
	return builder.String(), nil
}

func main() {
	message, err := processArgs()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	allCorrect := true
	var outMsg strings.Builder
	for _, word := range spellchecker.Check(message) {
		var txt aurora.Value
		if word.Correct {
			txt = aurora.Green(word.Text)
		} else {
			txt = aurora.Red(word.Text)
		}

		fmt.Fprintf(&outMsg, "%s ", txt)

		if !word.Correct {
			allCorrect = false
		}
	}

	if !allCorrect {
		fmt.Printf("Message: %s\n", outMsg.String())
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("\nMade a mistake commit anyway? y\\n? ")
		text, _ := reader.ReadString('\n')
		text = strings.ToLower(text[:len(text)-1])
		if text == "n" || text == "f" {
			return
		}
	}

	// parmas := strings.Split(fmt.Sprintf("git commit -m \"%s\"", message), " ")
	parmas := []string{"commit", "-m", message}
	stdout, err := exec.Command("git", parmas...).Output()

	if err != nil {
		fmt.Printf("Error %s\nMessage: %s\n", err.Error(), string(stdout))
		return
	}

	fmt.Printf("success %s\n", string(stdout))
}
