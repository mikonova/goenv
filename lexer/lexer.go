package lexer

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
)

var syntaxError error = errors.New("\033[0;31m[ERROR]\033[0;37m goenv: incorrect syntax: ")

func FetchString(file *os.File) {
	keymap := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Fatalln("\033[0;31m[ERROR]\033[0;37m file scanning error")
		}
		str := scanner.Text()
		str = strings.TrimSpace(str)
		runes := []rune(str)
		if !strings.Contains(str, "=") && runes[0] != '#' {
			log.Fatal(syntaxError, str, "\n")
		}
		if runes[0] == '=' {
			log.Fatal(syntaxError, string(runes), "\n")
		}
		if runes[0] != '#' {
			key, val := tokenize(runes)
			keymap[key] = val
		}

	}
	fillEnv(keymap)
}

func tokenize(runes []rune) (keyToken string, valToken string) {
	var (
		key, value, valueUnstripped string
	)
	for idx := range runes {

		if runes[idx] == '=' {
			key = string(runes[0:idx])
			key = strings.TrimSpace(key)
			valueStart := idx + 1
			valueUnstripped = string(runes[valueStart:])
			//fmt.Println(valueUnstripped)
			break
		}
	}
	// we have quotes
	if strings.ContainsAny(valueUnstripped, "\"'") {
		runeValue := []rune(valueUnstripped)[1:]
		if index, isFound := searchVal(runeValue, "\"'"); isFound {
			value = string(runeValue[:index])
		}
		// we dont have quotes
	} else {
		runeValue := []rune(valueUnstripped)
		if index, isFound := searchVal(runeValue, "#\n "); isFound {
			value = string(runeValue[:index])
		} else {
			value = string(runeValue)
		}
	}
	return key, value
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
