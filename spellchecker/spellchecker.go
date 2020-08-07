package spellchecker

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var (
	wordRe   *regexp.Regexp
	wordDict map[string]bool
)

//Word contains spelling info on a word
type Word struct {
	Text    string
	Correct bool
}

func init() {
	wordRe = regexp.MustCompile("^[A-Z]?[a-z]+$")
	wordDict = make(map[string]bool)

	file, err := os.Open("words_alpha.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wordDict[strings.ToLower(scanner.Text())] = true
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

}

func processWord(in string) Word {
	wLower := strings.ToLower(in)
	wInfo := Word{}
	wInfo.Text = in
	if wordRe.Match([]byte(in)) {
		_, ok := wordDict[wLower]
		wInfo.Correct = ok
	} else {
		wInfo.Correct = true
	}
	return wInfo
}

//Check returns an slice of Words with if it's marked as a mistake or not
func Check(msg string) []Word {

	splits := strings.Split(msg, " ")

	result := make([]Word, len(splits))

	for i, word := range splits {
		result[i] = processWord(word)
	}

	return result
}
