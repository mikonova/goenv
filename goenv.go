package goenv

import (
	"log"
	"os"

	"github.com/mikonova/goenv/lexer"
)

// Contains all file handles collected before ParseFiles () is called
var FileHandles []*os.File

// Collects all file relative paths into slice of handles, ready to work. Can be called multiple times
func FetchFiles(path ...string) {
	for index := range path {
		file, err := os.OpenFile(path[index], os.O_RDONLY, 0640)
		if err != nil {
			log.Fatalln("\033[0;31m[ERROR]\033[0;37m goenv: file not found")
		}
		FileHandles = append(FileHandles, file)
	}

}

// parses all variables from .env to the go environment variables
func ParseFiles() {
	for index := range FileHandles {
		lexer.FetchStrings(FileHandles[index])
	}
}
