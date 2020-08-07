package spellchecker

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
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

const dictFileURL = "https://raw.githubusercontent.com/dwyl/english-words/master/words_alpha.txt"

func downloadDictFile(filePath string) {
	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(dictFileURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("bad status: %s", resp.Status))
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}

func init() {
	wordRe = regexp.MustCompile("^[A-Z]?[a-z]+$")
	wordDict = make(map[string]bool)

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	appDataPath := filepath.Join(usr.HomeDir, ".scommit")
	dictFilePath := filepath.Join(appDataPath, "dict.txt")
	if _, err := os.Stat(appDataPath); os.IsNotExist(err) {
		fmt.Printf("First time Running!\nMaking program data folder at %s\n", appDataPath)
		err := os.Mkdir(appDataPath, os.FileMode(0777))
		if err != nil {
			panic(err)
		}
	}
	if _, err := os.Stat(dictFilePath); os.IsNotExist(err) {
		fmt.Printf("Downloading dict file\n")
		downloadDictFile(dictFilePath)
		fmt.Printf("dict file downloaded\n")
	}
	file, err := os.Open(dictFilePath)
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
