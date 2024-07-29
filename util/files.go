package util

import (
    "bufio"
    "io"
    "os"
    "path"
)

func MustLoad(open *os.File, err error) io.ReaderAt {
    if err != nil {
        panic(err)
    }
    return open
}
func MustOpen(filename string) io.ReadCloser {
    open, _ := os.Open(filename)
    return open
}

func FileExists(filename string) bool {
    _, err := os.Stat(filename)
    return !os.IsNotExist(err)
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
