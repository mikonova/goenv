package lexer

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
)

var syntaxError error = errors.New("\033[0;31m[ERROR]\033[0;37m goenv: incorrect syntax: ")

func FetchStrings(file *os.File) {
	keymap := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Fatalln("\033[0;31m[ERROR]\033[0;37m file scanning error")
		}
		str := scanner.Text()
		str = strings.TrimSpace(str)
		str += "\n"
		runes := []rune(str)
		if !strings.Contains(str, "=") && runes[0] != '#' {
			log.Fatal(syntaxError, str, "\n")
		}
		if runes[0] == '=' {
			log.Fatal(syntaxError, string(runes), "\n")
		}
		if runes[0] != '#' {
			key, val, err := tokenize(runes)
			if err != nil {
				log.Fatalln(err)
			}
			keymap[key] = val
		}

	}
	fillEnv(keymap)
}

func tokenize(runes []rune) (string, string, error) {
	var (
		key, value      string
		valueUnstripped []rune
	)
	if index, isValid := searchVal(runes, "="); isValid {
		key = string(runes[:index])
		valueUnstripped = runes[index+1:]

		if matchSlice, matches := findAll(valueUnstripped, "\"'"); matches == 2 {
			value = string(valueUnstripped[matchSlice[0]+1 : matchSlice[1]])
		} else if _, matches := findAll(valueUnstripped, "\"'"); matches == 1 {
			return "", "", errors.New(syntaxError.Error() + "cannot find closing quote in " + string(runes) + "\n")
		} else if idx, isEndl := searchVal(valueUnstripped, "#\n "); isEndl {
			value = string(valueUnstripped[:idx])
		} else {
			return "", "", errors.New(syntaxError.Error() + string(runes) + "\n")
		}
	}
	return key, value, nil
}

// searches for all inclusions of a rune in a slice of runes, returns list of indexes and number of iclusions
func findAll(source []rune, sample string) (matchIndexes []int, inclusions int) {
	sampleRuneSlice := []rune(sample)
	for k, v := range source {
		if findAllSearch(v, sampleRuneSlice) {
			inclusions++
			matchIndexes = append(matchIndexes, k)
		}
	}
	return matchIndexes, inclusions
}

// search function for findAll
func findAllSearch(char rune, compare []rune) bool {
	for _, val := range compare {
		if char == val {
			return true
		}
	}
	return false
}

// searches for the first occurrence of any char of sample in a slice of runes
func searchVal(source []rune, sample string) (int, bool) {
	for index := range source {
		if strings.ContainsAny(string(source[index]), sample) {
			return index, true
		}
	}
	return 0, false
}

func fillEnv(envList map[string]string) {
	for key, val := range envList {
		if err := os.Setenv(key, val); err != nil {
			log.Panic("\033[0;31m[ERROR]\033[0;37m goenv: failed to load variable to environment\n")
		}
	}
}
