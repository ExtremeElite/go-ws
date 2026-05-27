package util

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func PathToEveryOne(path string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal("path error")
	}
	baseDir := strings.Replace(dir, `\\`, "/", -1)
	p, _ := filepath.Abs(baseDir + "/" + path)
	return p
}
