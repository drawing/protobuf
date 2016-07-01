package main

import (
	"flag"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	var inPath string
	var outPath string

	flag.StringVar(&inPath, "in", ".", "input dir")
	flag.StringVar(&outPath, "out", ".", "input dir")
	flag.Parse()

	log.Println("start...")

	dir, err := os.Open(inPath)
	if err != nil {
		log.Fatalln("open path", inPath, "failed:", err)
	}
	defer dir.Close()

	names, err := dir.Readdirnames(0)
	if err != nil {
		log.Fatalln("read files name failed", err)
	}

	for _, name := range names {
		if !strings.HasSuffix(name, ".proto") {
			continue
		}

		goname := strings.TrimSuffix(name, ".proto") + ".pb.go"

		log.Println("work on:", name)

		pbpath := inPath + "/" + name
		tmppath := os.TempDir() + "/" + goname

		task := exec.Command("protoc", "--go_out="+os.TempDir(), pbpath)
		err := task.Run()
		if err != nil {
			log.Fatalln("exec protoc failed:", err)
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, tmppath, nil, parser.ImportsOnly)
		if err != nil {
			log.Fatalln("go file parse failed:", err)
		}

		pkgname := f.Name.Name

		log.Println("packagename=", pkgname)

		os.Mkdir(outPath+"/"+pkgname, os.ModePerm)

		os.Rename(tmppath, outPath+"/"+goname)
	}

	log.Println("end")
}
