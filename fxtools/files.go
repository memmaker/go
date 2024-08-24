package fxtools

import (
	"bufio"
	"io"
	"os"
	"path"
	"strings"
)

func MustLoad(open *os.File, err error) *os.File {
	if err != nil {
		panic(err)
	}
	return open
}
func MustOpen(filename string) *os.File {
	return MustLoad(os.Open(filename))
}
func OpenOrEmpty(filename string) *os.File {
	open, err := os.Open(filename)
	if err != nil { // return fake file
		return &os.File{}
	}
	return open
}

func MustCreate(filename string) *os.File {
	return MustLoad(os.Create(filename))
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func DirExists(dir string) bool {
	fi, err := os.Stat(dir)
	return err == nil && fi.IsDir()
}

func ReadFileAsLines(filename string) []string {
	var lines []string
	file, err := os.Open(filename)
	if err != nil {
		return lines
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
func ReadFile(filename string) string {
	var sb strings.Builder
	file, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
		sb.WriteString("\n")
	}
	return sb.String()
}

func CreateFile(file string) io.WriteCloser {
	f, err := os.Create(file)
	if err != nil {
		return nil
	}
	return f
}

func FilesInDirByExtension(dir string, extension string) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var result []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if extension != "" && extension != file.Name()[len(file.Name())-len(extension):] {
			continue
		}
		result = append(result, path.Join(dir, file.Name()))
	}
	return result

}
