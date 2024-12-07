package main

import (
	"github.com/memmaker/go/fxtools"
	"github.com/memmaker/go/recfile"
	"os"
	"path"
	"strings"
)

type MyRecord struct {
	Name      string
	Age       int
	Height    float64
	IsStudent bool
	Address   Address
}

type Address struct {
	Street string
	City   string
	Zip    int
}

func main() {
	pathname := "/Users/felix/Projects/Contractor/data_atom/definitions"

	repoman := recfile.NewRepoMan()

	// iterate over all files in the directory
	files, err := os.ReadDir(pathname)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		lowerName := strings.ToLower(file.Name())
		if !strings.HasSuffix(lowerName, ".rec") {
			continue
		}
		println("Processing", file.Name())
		filePath := path.Join(pathname, file.Name())
		openFile := fxtools.MustOpen(filePath)
		repoman.AddRepo(openFile)
	}

	println(repoman.Status())
}
