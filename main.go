package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var test bool
var verbose bool

func fixPerms(fullPath string, info os.FileInfo, err error) error {
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return err
	}

	contentInfoBuf, err := exec.Command("file", "-b", fullPath).Output()
	if err != nil {
		return fmt.Errorf("[%v, %v] file: %v", fullPath, info, err)
	}
	fileInfo := strings.TrimSpace(string(contentInfoBuf))

	var mode os.FileMode = 0644
	var prefix string = "default"
	const (
		elf64Executable = "ELF 64-bit LSB executable"
		elf32Executable = "ELF 32-bit LSB executable"
		asciiExecutable = "ASCII text executable"
		shebang         = "#!"
	)
	switch {
	case info.IsDir():
		mode = 0755 + os.ModeDir
		prefix = "dir"
	case strings.Contains(fileInfo, elf32Executable) || strings.Contains(fileInfo, elf64Executable):
		mode = 0755
		prefix = "binary-exec"
	case strings.Contains(fileInfo, asciiExecutable):
		f, err := os.Open(fullPath)
		if err != nil {
			return fmt.Errorf("[%v, %v] open file: %v", fullPath, info.Name(), err)
		}
		defer func() { _ = f.Close() }()
		buf := make([]byte, 2)
		if _, err := f.Read(buf); err == nil {
			if bytes.Equal(buf, []byte(shebang)) {
				mode = 0755
				prefix = "shebang-exec"
			}
		}
	}
	if !test {
		if err := os.Chmod(fullPath, mode); err != nil {
			return fmt.Errorf("[%v, %v] %v chmod %v: %v", fullPath, info.Name(), prefix, mode, err)
		}
	}
	if (test || verbose) && info.Mode() != mode {
		fmt.Printf("[%v %v -> %v] %v :: %v\n", prefix, info.Mode(), mode, fileInfo, fullPath)
	}
	return nil
}

func main() {
	var root string
	flag.StringVar(&root, "root", "", "Target folder to reset file permissions to default.")
	flag.BoolVar(&verbose, "verbose", false, "Display performed changes.")
	flag.BoolVar(&test, "test", false, "Display planned changes, but don't do change.")
	flag.Parse()

	if root == "" {
		os.Stderr.WriteString("'root' argument is mandatory to specify.")
		os.Exit(1)
	}

	if err := filepath.Walk(root, fixPerms); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
