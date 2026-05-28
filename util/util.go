package util

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ResolvePath(relativePath string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal("failed to resolve working directory: ", err)
	}
	baseDir := strings.Replace(dir, `\\`, "/", -1)
	p, _ := filepath.Abs(baseDir + "/" + relativePath)
	return p
}
